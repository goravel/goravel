package modify

import (
	"go/token"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"

	contractsmatch "github.com/goravel/framework/contracts/packages/match"
	"github.com/goravel/framework/contracts/packages/modify"
	"github.com/goravel/framework/packages/match"
	supportfile "github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
)

// withSliceConfig holds configuration for adding items to a slice in foundation.Setup() chain.
type withSliceConfig struct {
	// fileName is the name of the file to create (e.g., "commands.go", "migrations.go")
	fileName string
	// withMethodName is the name of the method in foundation.Setup() chain (e.g., "WithCommands", "WithMigrations")
	withMethodName string
	// helperFuncName is the name of the helper function (e.g., "Commands", "Migrations")
	helperFuncName string
	// typePackage is the package name for the slice element type (e.g., "console", "schema")
	typePackage string
	// typeName is the type name for the slice element (e.g., "Command", "Migration")
	typeName string
	// typeImportPath is the import path for the type package (e.g., "github.com/goravel/framework/contracts/console")
	typeImportPath string
	// fileExistsError is the error to return when the file exists but WithMethod is not registered
	fileExistsError error
	// stubTemplate is the function that returns the stub content for creating the file
	stubTemplate func() string
	// matcherFunc is the function that returns the matcher to find the target location (e.g., matchPkg.Commands, matchPkg.Migrations)
	matcherFunc func() []contractsmatch.GoNode
}

// withSliceHandler handles adding items to a slice in foundation.Setup() chain.
// It supports two modes:
//  1. Helper file mode (alwaysInline=false): Items can be managed in separate helper files (e.g., commands.go, migrations.go)
//  2. Always inline mode (alwaysInline=true): Items are always added directly to app.go
type withSliceHandler struct {
	config      withSliceConfig
	appFilePath string
	filePath    string
	fileExists  bool
}

// newWithSliceHandler creates a new withSliceHandler with the given configuration.
func newWithSliceHandler(config withSliceConfig) *withSliceHandler {
	appFilePath := path.Bootstrap("app.go")
	bootstrapDir := filepath.Dir(appFilePath)

	var filePath string
	var fileExists bool

	// Only construct file path and check existence for helper file mode
	if config.fileName != "" {
		filePath = filepath.Join(bootstrapDir, config.fileName)
		fileExists = supportfile.Exists(filePath)
	}

	return &withSliceHandler{
		config:      config,
		appFilePath: appFilePath,
		filePath:    filePath,
		fileExists:  fileExists,
	}
}

// AddItem adds an item to a slice in the foundation.Setup() chain.
//
// Behavior depends on the configuration mode:
//
// Always-inline mode (alwaysInline=true):
//   - Always adds items directly to the inline array in app.go
//   - Never creates helper files
//   - Used for routing (func{})
//
// Helper file mode (alwaysInline=false):
//  1. If WithMethod doesn't exist:
//     - If helper file exists: Returns an error (file should only exist when WithMethod is registered)
//     - If helper file doesn't exist: Creates the helper file with the item, adds WithMethod(HelperFunc) to Setup()
//  2. If WithMethod exists:
//     - If helper file exists: Appends the item to the helper function
//     - If helper file doesn't exist: Appends the item to the inline array in app.go
//
// Parameters:
//   - pkg: Package path of the item (e.g., "github.com/goravel/app/rules" or "goravel/routes")
//   - item: Item expression to add (e.g., "&rules.Uppercase{}" or "routes.Web")
//
// Example 1 - Creating WithMethod with helper file:
//
// Before (app.go):
//
//	foundation.Setup().WithConfig(config.Boot).Start()
//
// After (app.go):
//
//	foundation.Setup().WithRules(Rules).WithConfig(config.Boot).Start()
//
// And creates helper file (e.g., bootstrap/rules.go):
//
//	package bootstrap
//	import "github.com/goravel/framework/contracts/validation"
//	func Rules() []validation.Rule {
//	  return []validation.Rule{&rules.Uppercase{}}
//	}
//
// Example 2 - Appending to inline array:
//
// Before (app.go):
//
//	foundation.Setup().WithRules(func() []validation.Rule{
//	  return []validation.Rule{
//		&rules.ExistingRule{},
//	  },
//	}).Start()
//
// After (app.go):
//
//	foundation.Setup().WithRules(func() []validation.Rule{
//	  return []validation.Rule{
//		&rules.ExistingRule{},
//		&rules.Uppercase{},
//	  },
//	}).Start()
//
// Example 3 - Appending to helper function:
//
// If helper file exists with Rules() function, appends to that function instead.
func (r *withSliceHandler) AddItem(pkg, item string) error {
	hasFoundationSetup, err := r.containsFoundationSetupInAppFile()
	if err != nil {
		return err
	}
	if !hasFoundationSetup {
		return nil
	}

	if r.config.fileName != "" {
		// Helper file mode (like commands/migrations)
		withMethodExists, err := r.checkWithMethodExists()
		if err != nil {
			return err
		}

		if !withMethodExists {
			if r.fileExists {
				return r.config.fileExistsError
			}

			if err := r.createFile(); err != nil {
				return err
			}

			if err := r.addItemToFile(pkg, item); err != nil {
				return err
			}

			return GoFile(r.appFilePath).Find(match.FoundationSetup()).Modify(r.setupWithFunction()).Apply()
		}

		if r.fileExists {
			if err := r.addItemToFile(pkg, item); err != nil {
				return err
			}
			return nil
		}
	}

	if err := r.addImports(pkg); err != nil {
		return err
	}

	return GoFile(r.appFilePath).Find(match.FoundationSetup()).Modify(r.setupInline(item)).Apply()
}

// RemoveItem removes an item from the slice in foundation.Setup() chain.
//
// In always-inline mode: Removes the item from the inline array in app.go and cleans up imports.
// In helper file mode:
//   - If the helper file exists (e.g., providers.go), it removes the item from that file
//   - If the helper file doesn't exist, it removes the item from the inline array in app.go
//   - Cleans up unused imports after removing the item
func (r *withSliceHandler) RemoveItem(pkg, item string) error {
	// Check if foundation.Setup() exists in app.go before performing any actions
	hasFoundationSetup, err := r.containsFoundationSetupInAppFile()
	if err != nil {
		return err
	}
	if !hasFoundationSetup {
		return nil
	}

	withMethodExists, err := r.checkWithMethodExists()
	if err != nil {
		return err
	}

	if !withMethodExists {
		// If WithMethod doesn't exist, there's nothing to remove
		return nil
	}

	// Helper file mode (like commands/migrations)
	if r.fileExists {
		// Remove from helper file
		if err := r.removeItemFromFile(pkg, item); err != nil {
			return err
		}
		return nil
	}

	// Remove from inline array in app.go
	if err := GoFile(r.appFilePath).Find(match.FoundationSetup()).Modify(r.removeInline(item)).Apply(); err != nil {
		return err
	}

	// Clean up imports
	return r.removeImports(pkg)
}

// checkWithMethodExists checks if the WithMethod exists in the foundation.Setup() chain.
//
// Example: For a config with withMethodName="WithCommands", it searches for:
//
//	foundation.Setup().WithCommands(Commands).Boot()
//
// Returns true if ".WithCommands" is found in the app file content.
func (r *withSliceHandler) checkWithMethodExists() (bool, error) {
	content, err := supportfile.GetContent(r.appFilePath)
	if err != nil {
		return false, err
	}

	return strings.Contains(content, r.config.withMethodName), nil
}

// containsFoundationSetupInAppFile checks if .Setup(). exists in the app file.
//
// This check is performed at the start of AddItem to ensure no actions are taken
// if the .Setup(). call doesn't exist in the file.
//
// Returns true if ".Setup()." is found in the app file content.
func (r *withSliceHandler) containsFoundationSetupInAppFile() (bool, error) {
	content, err := supportfile.GetContent(r.appFilePath)
	if err != nil {
		return false, err
	}

	return strings.Contains(content, ".Setup()."), nil
}

// createFile creates the file with the helper function.
//
// Example: Creates a new file like "bootstrap/commands.go" with content:
//
//	func Commands() []console.Command {
//	    return []console.Command{}
//	}
//
// The actual content is generated by config.stubTemplate().
func (r *withSliceHandler) createFile() error {
	return supportfile.PutContent(r.filePath, r.config.stubTemplate())
}

// addImports adds the required imports for the item package and type package.
//
// Example: For pkg="github.com/user/app/commands" and typeImportPath="github.com/goravel/framework/contracts/console":
//
//	import (
//	    "github.com/user/app/commands"
//	    "github.com/goravel/framework/contracts/console"
//	)
func (r *withSliceHandler) addImports(pkg string) error {
	importMatchers := match.Imports()
	if err := GoFile(r.appFilePath).FindOrCreate(importMatchers, createImport).Modify(AddImport(pkg)).Apply(); err != nil {
		return err
	}

	// Skip adding type import for function slices (like routing)
	if r.config.typeImportPath != "" {
		return GoFile(r.appFilePath).Find(importMatchers).Modify(AddImport(r.config.typeImportPath)).Apply()
	}

	return nil
}

// addItemToFile adds an item to the existing helper function in the file.
//
// Example: For pkg="github.com/user/app/commands" and item="&commands.MyCommand{}":
//
// Before:
//
//	func Commands() []console.Command {
//	    return []console.Command{}
//	}
//
// After:
//
//	import "github.com/user/app/commands"
//
//	func Commands() []console.Command {
//	    return []console.Command{
//	        &commands.MyCommand{},
//	    }
//	}
func (r *withSliceHandler) addItemToFile(pkg, item string) error {
	// Add the item package import
	importMatchers := match.Imports()
	if err := GoFile(r.filePath).FindOrCreate(importMatchers, createImport).Modify(AddImport(pkg)).Apply(); err != nil {
		return err
	}

	// Add the item to the helper function using the provided matcher function
	return GoFile(r.filePath).Find(r.config.matcherFunc()).Modify(Register(item)).Apply()
}

// removeImports removes the item package import if it's no longer used.
// It checks both app.go and the helper file (if it exists) to determine if the import is still in use.
func (r *withSliceHandler) removeImports(pkg string) error {
	importMatchers := match.Imports()
	return GoFile(r.appFilePath).Find(importMatchers).Modify(RemoveImport(pkg)).Apply()
}

// removeItemFromFile removes an item from the existing helper function in the file.
//
// Example: For pkg="github.com/user/app/commands" and item="&commands.MyCommand{}":
//
// Before:
//
//	import "github.com/user/app/commands"
//
//	func Commands() []console.Command {
//	    return []console.Command{
//	        &commands.MyCommand{},
//	        &commands.OtherCommand{},
//	    }
//	}
//
// After:
//
//	import "github.com/user/app/commands"
//
//	func Commands() []console.Command {
//	    return []console.Command{
//	        &commands.OtherCommand{},
//	    }
//	}
func (r *withSliceHandler) removeItemFromFile(pkg, item string) error {
	// Remove the item from the helper function using the provided matcher function
	if err := GoFile(r.filePath).Find(r.config.matcherFunc()).Modify(Unregister(item)).Apply(); err != nil {
		return err
	}

	// Clean up the import if it's no longer used
	importMatchers := match.Imports()
	return GoFile(r.filePath).Find(importMatchers).Modify(RemoveImport(pkg)).Apply()
}

// appendToExisting appends an item to an existing WithMethod call.
//
// Example: For withCall representing ".WithCommands(func() []console.Command{...})" and itemExpr="&commands.Cmd2{}":
//
// Before:
//
//	.WithCommands(func() []console.Command{
//	    return []console.Command{
//	        &commands.Cmd1{},
//	    }
//	})
//
// After:
//
//	.WithCommands(func() []console.Command{
//	    return []console.Command{
//	        &commands.Cmd1{},
//	        &commands.Cmd2{},
//	    }
//	})
func (r *withSliceHandler) appendToExisting(withCall *dst.CallExpr, itemExpr dst.Expr) {
	if len(withCall.Args) == 0 {
		return
	}

	// Check if the argument is a function literal (new API: func() []Type)
	funcLit, ok := withCall.Args[0].(*dst.FuncLit)
	if !ok {
		return
	}

	// Find the return statement and its composite literal
	for _, stmt := range funcLit.Body.List {
		if retStmt, ok := stmt.(*dst.ReturnStmt); ok {
			if len(retStmt.Results) > 0 {
				if compositeLit, ok := retStmt.Results[0].(*dst.CompositeLit); ok {
					// Add proper formatting for multi-line arrays
					if len(compositeLit.Elts) > 0 {
						compositeLit.Elts[0].Decorations().Before = dst.NewLine
					}

					// Add newline before new item and after (for closing brace)
					itemExpr.Decorations().Before = dst.NewLine
					itemExpr.Decorations().After = dst.NewLine

					compositeLit.Elts = append(compositeLit.Elts, itemExpr)
					return
				}
			}
		}
	}
}

// createWithMethod creates a new WithMethod call and inserts it into the chain.
// Wraps the composite literal in a function literal: func() []Type { return []Type{...} }
//
// Before:
//
//	foundation.Setup().Boot()
//
// After:
//
//	foundation.Setup().
//	    WithCommands(func() []console.Command {
//	        return []console.Command{
//	            &commands.MyCommand{},
//	        }
//	    }).
//	    Boot()

func (r *withSliceHandler) createWithMethod(setupCall *dst.CallExpr, parentOfSetup *dst.SelectorExpr, itemExpr dst.Expr) {
	var arrayType dst.Expr
	var returnArrayType dst.Expr

	// For typed slices (commands, migrations, etc.)
	arrayType = &dst.SelectorExpr{
		X:   &dst.Ident{Name: r.config.typePackage},
		Sel: &dst.Ident{Name: r.config.typeName},
	}
	// Create a separate instance for the return type
	returnArrayType = &dst.SelectorExpr{
		X:   &dst.Ident{Name: r.config.typePackage},
		Sel: &dst.Ident{Name: r.config.typeName},
	}

	// Add proper formatting for multi-line array
	itemExpr.Decorations().Before = dst.NewLine
	itemExpr.Decorations().After = dst.NewLine

	compositeLit := &dst.CompositeLit{
		Type: &dst.ArrayType{
			Elt: arrayType,
		},
		Elts: []dst.Expr{itemExpr},
	}

	// For typed slices, wrap the composite literal in a function literal: func() []Type { return []Type{...} }
	funcArg := &dst.FuncLit{
		Type: &dst.FuncType{
			Params: &dst.FieldList{},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Type: &dst.ArrayType{
							Elt: returnArrayType,
						},
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{compositeLit},
					Decs: dst.ReturnStmtDecorations{
						NodeDecs: dst.NodeDecs{
							Before: dst.NewLine,
							After:  dst.NewLine,
						},
					},
				},
			},
		},
	}

	newWithCall := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X: setupCall,
			Sel: &dst.Ident{
				Name: r.config.withMethodName,
				Decs: dst.IdentDecorations{
					NodeDecs: dst.NodeDecs{
						Before: dst.NewLine,
					},
				},
			},
		},
		Args: []dst.Expr{funcArg},
	}

	parentOfSetup.X = newWithCall
}

// findFoundationSetupCalls walks the chain to find Setup() and WithMethod() calls.
//
// Example: For a chain like:
//
//	foundation.Setup().WithCommands(Commands).WithMigrations(Migrations).Boot()
//
// If withMethodName="WithCommands", returns:
//   - setupCall: the CallExpr for "foundation.Setup()"
//   - withCall: the CallExpr for "WithCommands(Commands)" (if it exists)
//   - parentOfSetup: the SelectorExpr representing the next method after Setup (e.g., ".WithCommands")
func (r *withSliceHandler) findFoundationSetupCalls(callExpr *dst.CallExpr) (setupCall, withCall *dst.CallExpr, parentOfSetup *dst.SelectorExpr) {
	current := callExpr
	for current != nil {
		if sel, ok := current.Fun.(*dst.SelectorExpr); ok {
			if innerCall, ok := sel.X.(*dst.CallExpr); ok {
				if innerSel, ok := innerCall.Fun.(*dst.SelectorExpr); ok {
					// Check if this is the Setup() call
					if innerSel.Sel.Name == "Setup" {
						if ident, ok := innerSel.X.(*dst.Ident); ok && ident.Name == "foundation" {
							setupCall = innerCall
							parentOfSetup = sel
							break
						}
					}
					// Check if this is the WithMethod
					if innerSel.Sel.Name == r.config.withMethodName {
						withCall = innerCall
					}
				}
				current = innerCall
				continue
			}
		}
		break
	}
	return
}

// setupInline returns an action that modifies the foundation.Setup() chain (inline array).
//
// This is used in two scenarios:
//  1. Always-inline mode: Always adds items directly to the inline array (e.g., routing)
//  2. Helper file mode: When WithMethod exists in app.go but the helper file doesn't exist
//
// Supports both typed slices ([]console.Command) and function slices ([]func{}).
//
// Example 1 - Typed slice (commands) with existing WithCommands:
//
// Before:
//
//		foundation.Setup().
//		    WithCommands(func []console.Command {
//		        return []console.Command{
//		        	&commands.ExistingCmd{},
//	            }
//		    }).
//		    Boot()
//
// After:
//
//		foundation.Setup().
//		    WithCommands(func []console.Command {
//		        return []console.Command{
//		        	&commands.ExistingCmd{},
//		        	&commands.MyCommand{},
//	           }
//		    }).
//		    Boot()
//
// Or if WithCommands doesn't exist:
//
// Before:
//
//	foundation.Setup().Boot()
//
// After:
//
//	foundation.Setup().
//	    WithCommands(func []console.Command {
//	        return []console.Command{
//	        	&commands.MyCommand{},
//	        }
//	    }).
//	    Boot()
//
// Example 2 - Function slice (routing):
//
// Before:
//
//	foundation.Setup().Boot()
//
// After:
//
//	foundation.Setup().
//	    WithRouting(func(){
//	        routes.Web(),
//	    }).
//	    Boot()
func (r *withSliceHandler) setupInline(item string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var callExpr *dst.CallExpr

		switch stmt := cursor.Node().(type) {
		case *dst.ExprStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			var ok bool
			callExpr, ok = stmt.X.(*dst.CallExpr)
			if !ok {
				return
			}
		case *dst.ReturnStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			if len(stmt.Results) == 0 {
				return
			}
			var ok bool
			callExpr, ok = stmt.Results[0].(*dst.CallExpr)
			if !ok {
				return
			}
		default:
			return
		}

		setupCall, withCall, parentOfSetup := r.findFoundationSetupCalls(callExpr)
		if setupCall == nil || parentOfSetup == nil {
			return
		}

		itemExpr := MustParseExpr(item).(dst.Expr)

		if withCall != nil {
			r.appendToExisting(withCall, itemExpr)
		} else {
			r.createWithMethod(setupCall, parentOfSetup, itemExpr)
		}
	}
}

// setupWithFunction returns an action that adds WithMethod(HelperFunc) to the foundation.Setup() chain.
//
// This is used when the helper file is created and we need to wire it into app.go using a helper function reference.
//
// Example: For withMethodName="WithCommands" and helperFuncName="Commands":
//
// Before:
//
//	foundation.Setup().Boot()
//
// After:
//
//	foundation.Setup().
//	    WithCommands(Commands).
//	    Boot()
//
// Where Commands is a helper function defined in bootstrap/commands.go that returns []console.Command.
func (r *withSliceHandler) setupWithFunction() modify.Action {
	return func(cursor *dstutil.Cursor) {
		var callExpr *dst.CallExpr

		switch stmt := cursor.Node().(type) {
		case *dst.ExprStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			var ok bool
			callExpr, ok = stmt.X.(*dst.CallExpr)
			if !ok {
				return
			}
		case *dst.ReturnStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			if len(stmt.Results) == 0 {
				return
			}
			var ok bool
			callExpr, ok = stmt.Results[0].(*dst.CallExpr)
			if !ok {
				return
			}
		default:
			return
		}

		setupCall, _, parentOfSetup := r.findFoundationSetupCalls(callExpr)
		if setupCall == nil || parentOfSetup == nil {
			return
		}

		// Create WithMethod(HelperFunc) call
		newWithCall := &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X: setupCall,
				Sel: &dst.Ident{
					Name: r.config.withMethodName,
					Decs: dst.IdentDecorations{
						NodeDecs: dst.NodeDecs{
							Before: dst.NewLine,
						},
					},
				},
			},
			Args: []dst.Expr{
				&dst.Ident{Name: r.config.helperFuncName},
			},
		}

		// Insert WithMethod into the chain
		parentOfSetup.X = newWithCall
	}
}

// removeInline returns an action that removes an item from the inline array in the foundation.Setup() chain.
//
// This is used when the helper file doesn't exist and items are stored directly in app.go.
// It removes the specified item from the inline array.
//
// Example: For item="&commands.MyCommand{}" with existing WithCommands:
//
// Before:
//
//	foundation.Setup().
//	    WithCommands([]console.Command{
//	        &commands.ExistingCmd{},
//	        &commands.MyCommand{},
//	    }).
//	    Boot()
//
// After:
//
//	foundation.Setup().
//	    WithCommands([]console.Command{
//	        &commands.ExistingCmd{},
//	    }).
//	    Boot()
func (r *withSliceHandler) removeInline(item string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var callExpr *dst.CallExpr

		switch stmt := cursor.Node().(type) {
		case *dst.ExprStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			var ok bool
			callExpr, ok = stmt.X.(*dst.CallExpr)
			if !ok {
				return
			}
		case *dst.ReturnStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			if len(stmt.Results) == 0 {
				return
			}
			var ok bool
			callExpr, ok = stmt.Results[0].(*dst.CallExpr)
			if !ok {
				return
			}
		default:
			return
		}

		_, withCall, _ := r.findFoundationSetupCalls(callExpr)
		if withCall == nil {
			return
		}

		itemExpr := MustParseExpr(item).(dst.Expr)
		r.removeFromExisting(withCall, itemExpr)
	}
}

// removeFromExisting removes an item from an existing WithMethod call.
//
// Example: For withCall representing ".WithCommands(func() []console.Command{...})" and itemExpr="&commands.Cmd1{}":
//
// Before:
//
//	.WithCommands(func() []console.Command {
//	    return []console.Command{
//	        &commands.Cmd1{},
//	        &commands.Cmd2{},
//	    }
//	})
//
// After:
//
//	.WithCommands(func() []console.Command {
//	    return []console.Command{
//	        &commands.Cmd2{},
//	    }
//	})
func (r *withSliceHandler) removeFromExisting(withCall *dst.CallExpr, itemExpr dst.Expr) {
	if len(withCall.Args) == 0 {
		return
	}

	// Check if the argument is a function literal (new API: func() []Type)
	funcLit, ok := withCall.Args[0].(*dst.FuncLit)
	if !ok {
		return
	}

	// Find the return statement and its composite literal
	for _, stmt := range funcLit.Body.List {
		if retStmt, ok := stmt.(*dst.ReturnStmt); ok {
			if len(retStmt.Results) > 0 {
				if compositeLit, ok := retStmt.Results[0].(*dst.CompositeLit); ok {
					// Use slices.DeleteFunc to remove the matching item
					compositeLit.Elts = slices.DeleteFunc(compositeLit.Elts, func(ex dst.Expr) bool {
						return match.EqualNode(itemExpr).MatchNode(ex)
					})
					return
				}
			}
		}
	}
}

// addMiddlewareAppendCall adds a new handler.Append() call to the function literal.
func addMiddlewareAppendCall(funcLit *dst.FuncLit, middlewareArg dst.Expr) {
	// Add newline decorations to middleware argument for proper formatting
	middlewareArg.Decorations().Before = dst.NewLine
	middlewareArg.Decorations().After = dst.NewLine

	appendStmt := &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X:   &dst.Ident{Name: "handler"},
				Sel: &dst.Ident{Name: "Append"},
			},
			Args: []dst.Expr{middlewareArg},
			Decs: dst.CallExprDecorations{
				NodeDecs: dst.NodeDecs{
					Before: dst.NewLine,
					After:  dst.NewLine,
				},
			},
		},
	}
	funcLit.Body.List = append(funcLit.Body.List, appendStmt)
}

// addMiddlewareImports adds the required imports for middleware and configuration packages.
func addMiddlewareImports(appFilePath, pkg string) error {
	importMatchers := match.Imports()
	if err := GoFile(appFilePath).FindOrCreate(importMatchers, createImport).Modify(AddImport(pkg)).Apply(); err != nil {
		return err
	}

	configImportPath := "github.com/goravel/framework/contracts/foundation/configuration"
	return GoFile(appFilePath).Find(importMatchers).Modify(AddImport(configImportPath)).Apply()
}

// appendToExistingMiddleware appends middleware to an existing WithMiddleware call.
func appendToExistingMiddleware(withMiddlewareCall *dst.CallExpr, middlewareExpr dst.Expr) {
	if len(withMiddlewareCall.Args) == 0 {
		return
	}

	funcLit, ok := withMiddlewareCall.Args[0].(*dst.FuncLit)
	if !ok {
		return
	}

	appendCall := findMiddlewareAppendCall(funcLit)
	if appendCall != nil {
		// Ensure the first existing argument doesn't have a newline before it
		if len(appendCall.Args) > 0 {
			appendCall.Args[0].Decorations().Before = dst.None
		}

		// Add newline decorations to the new middleware for proper formatting
		middlewareExpr.Decorations().Before = dst.NewLine
		middlewareExpr.Decorations().After = dst.NewLine

		appendCall.Args = append(appendCall.Args, middlewareExpr)
	} else {
		addMiddlewareAppendCall(funcLit, middlewareExpr)
	}
}

// containsFoundationSetup checks if the given AST node contains a foundation.Setup() call.
// It works with any dst.Node, including *dst.ExprStmt and *dst.ReturnStmt.
func containsFoundationSetup(node dst.Node) bool {
	var foundSetup bool
	dst.Inspect(node, func(n dst.Node) bool {
		if call, ok := n.(*dst.CallExpr); ok {
			if sel, ok := call.Fun.(*dst.SelectorExpr); ok {
				if ident, ok := sel.X.(*dst.Ident); ok {
					if ident.Name == "foundation" && sel.Sel.Name == "Setup" {
						foundSetup = true
						return false
					}
				}
			}
		}
		return true
	})
	return foundSetup
}

func createImport(node dst.Node) error {
	importDecl := &dst.GenDecl{
		Tok: token.IMPORT,
	}

	f := node.(*dst.File)

	newDecls := make([]dst.Decl, 0, len(f.Decls)+1)
	newDecls = append(newDecls, f.Decls[0], importDecl) // package and import

	if len(f.Decls) > 1 {
		newDecls = append(newDecls, f.Decls[1:]...) // others
	}

	f.Decls = newDecls

	return nil
}

// createWithMiddleware creates a new WithMiddleware call and inserts it into the chain.
func createWithMiddleware(setupCall *dst.CallExpr, parentOfSetup *dst.SelectorExpr, middlewareExpr dst.Expr) {
	// Add newline decorations to middleware argument for proper formatting
	middlewareExpr.Decorations().Before = dst.NewLine
	middlewareExpr.Decorations().After = dst.NewLine

	funcLit := &dst.FuncLit{
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{{Name: "handler"}},
						Type: &dst.SelectorExpr{
							X:   &dst.Ident{Name: "configuration"},
							Sel: &dst.Ident{Name: "Middleware"},
						},
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   &dst.Ident{Name: "handler"},
							Sel: &dst.Ident{Name: "Append"},
						},
						Args: []dst.Expr{middlewareExpr},
						Decs: dst.CallExprDecorations{
							NodeDecs: dst.NodeDecs{
								Before: dst.NewLine,
								After:  dst.NewLine,
							},
						},
					},
				},
			},
		},
	}

	newWithMiddlewareCall := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X: setupCall,
			Sel: &dst.Ident{
				Name: "WithMiddleware",
				Decs: dst.IdentDecorations{
					NodeDecs: dst.NodeDecs{
						Before: dst.NewLine,
					},
				},
			},
		},
		Args: []dst.Expr{funcLit},
	}

	// Insert WithMiddleware into the chain
	parentOfSetup.X = newWithMiddlewareCall
}

// findFoundationSetupCallsForMiddleware walks the chain to find Setup() and WithMiddleware() calls.
func findFoundationSetupCallsForMiddleware(callExpr *dst.CallExpr) (setupCall, withMiddlewareCall *dst.CallExpr, parentOfSetup *dst.SelectorExpr) {
	current := callExpr
	for current != nil {
		if sel, ok := current.Fun.(*dst.SelectorExpr); ok {
			if innerCall, ok := sel.X.(*dst.CallExpr); ok {
				if innerSel, ok := innerCall.Fun.(*dst.SelectorExpr); ok {
					// Check if this is the Setup() call
					if innerSel.Sel.Name == "Setup" {
						if ident, ok := innerSel.X.(*dst.Ident); ok && ident.Name == "foundation" {
							setupCall = innerCall
							parentOfSetup = sel
							break
						}
					}
					// Check if this is WithMiddleware
					if innerSel.Sel.Name == "WithMiddleware" {
						withMiddlewareCall = innerCall
					}
				}
				current = innerCall
				continue
			}
		}
		break
	}
	return
}

// findMiddlewareAppendCall finds the handler.Append() call in the function literal.
func findMiddlewareAppendCall(funcLit *dst.FuncLit) *dst.CallExpr {
	for _, stmt := range funcLit.Body.List {
		if exprStmt, ok := stmt.(*dst.ExprStmt); ok {
			if call, ok := exprStmt.X.(*dst.CallExpr); ok {
				if sel, ok := call.Fun.(*dst.SelectorExpr); ok {
					if sel.Sel.Name == "Append" {
						return call
					}
				}
			}
		}
	}
	return nil
}

// foundationSetupMiddleware returns an action that modifies the foundation.Setup() chain.
func foundationSetupMiddleware(middleware string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var callExpr *dst.CallExpr

		switch stmt := cursor.Node().(type) {
		case *dst.ExprStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			var ok bool
			callExpr, ok = stmt.X.(*dst.CallExpr)
			if !ok {
				return
			}
		case *dst.ReturnStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			if len(stmt.Results) == 0 {
				return
			}
			var ok bool
			callExpr, ok = stmt.Results[0].(*dst.CallExpr)
			if !ok {
				return
			}
		default:
			return
		}

		setupCall, withMiddlewareCall, parentOfSetup := findFoundationSetupCallsForMiddleware(callExpr)
		if setupCall == nil || parentOfSetup == nil {
			return
		}

		middlewareExpr := MustParseExpr(middleware).(dst.Expr)

		if withMiddlewareCall != nil {
			appendToExistingMiddleware(withMiddlewareCall, middlewareExpr)
		} else {
			createWithMiddleware(setupCall, parentOfSetup, middlewareExpr)
		}
	}
}

// addRouteImports adds the required imports for the route package.
func addRouteImports(appFilePath, pkg string) error {
	importMatchers := match.Imports()
	return GoFile(appFilePath).FindOrCreate(importMatchers, createImport).Modify(AddImport(pkg)).Apply()
}

// foundationSetupRouting returns an action that modifies the foundation.Setup() chain to add or update WithRouting.
func foundationSetupRouting(route string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var callExpr *dst.CallExpr

		switch stmt := cursor.Node().(type) {
		case *dst.ExprStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			var ok bool
			callExpr, ok = stmt.X.(*dst.CallExpr)
			if !ok {
				return
			}
		case *dst.ReturnStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			if len(stmt.Results) == 0 {
				return
			}
			var ok bool
			callExpr, ok = stmt.Results[0].(*dst.CallExpr)
			if !ok {
				return
			}
		default:
			return
		}

		setupCall, withRoutingCall, parentOfSetup := findFoundationSetupCallsForRouting(callExpr)
		if setupCall == nil || parentOfSetup == nil {
			return
		}

		routeExpr := MustParseExpr(route).(dst.Expr)
		routeStmt := &dst.ExprStmt{X: routeExpr}

		if withRoutingCall != nil {
			appendToExistingRouting(withRoutingCall, routeStmt)
		} else {
			createWithRouting(setupCall, parentOfSetup, routeStmt)
		}
	}
}

// removeRouteFromSetup returns an action that removes a route from WithRouting in foundation.Setup().
func removeRouteFromSetup(route string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var callExpr *dst.CallExpr

		switch stmt := cursor.Node().(type) {
		case *dst.ExprStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			var ok bool
			callExpr, ok = stmt.X.(*dst.CallExpr)
			if !ok {
				return
			}
		case *dst.ReturnStmt:
			if !containsFoundationSetup(stmt) {
				return
			}
			if len(stmt.Results) == 0 {
				return
			}
			var ok bool
			callExpr, ok = stmt.Results[0].(*dst.CallExpr)
			if !ok {
				return
			}
		default:
			return
		}

		_, withRoutingCall, _ := findFoundationSetupCallsForRouting(callExpr)
		if withRoutingCall == nil || len(withRoutingCall.Args) == 0 {
			return
		}

		funcLit, ok := withRoutingCall.Args[0].(*dst.FuncLit)
		if !ok {
			return
		}

		routeExpr := MustParseExpr(route).(dst.Expr)
		routeStmt := &dst.ExprStmt{X: routeExpr}
		removeRouteFromFuncBody(funcLit.Body, routeStmt)
	}
}

// findFoundationSetupCallsForRouting finds the Setup() call and WithRouting() call in the chain.
func findFoundationSetupCallsForRouting(callExpr *dst.CallExpr) (setupCall *dst.CallExpr, withRoutingCall *dst.CallExpr, parentOfSetup *dst.SelectorExpr) {
	current := callExpr
	for current != nil {
		if sel, ok := current.Fun.(*dst.SelectorExpr); ok {
			if innerCall, ok := sel.X.(*dst.CallExpr); ok {
				if innerSel, ok := innerCall.Fun.(*dst.SelectorExpr); ok {
					// Check if this is the Setup() call
					if innerSel.Sel.Name == "Setup" {
						if ident, ok := innerSel.X.(*dst.Ident); ok && ident.Name == "foundation" {
							setupCall = innerCall
							parentOfSetup = sel
							break
						}
					}
					// Check if this is WithRouting
					if innerSel.Sel.Name == "WithRouting" {
						withRoutingCall = innerCall
					}
				}
				current = innerCall
				continue
			}
		}
		break
	}
	return
}

// createWithRouting creates a new WithRouting call and inserts it into the chain.
func createWithRouting(setupCall *dst.CallExpr, parentOfSetup *dst.SelectorExpr, routeStmt *dst.ExprStmt) {
	// Add proper decorations to the route statement for formatting
	routeStmt.Decs.Before = dst.NewLine
	routeStmt.Decs.After = dst.NewLine

	funcLit := &dst.FuncLit{
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{routeStmt},
		},
	}

	withRoutingCall := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X: setupCall,
			Sel: &dst.Ident{
				Name: "WithRouting",
				Decs: dst.IdentDecorations{
					NodeDecs: dst.NodeDecs{
						Before: dst.NewLine,
					},
				},
			},
		},
		Args: []dst.Expr{funcLit},
	}

	// Insert WithRouting into the chain
	parentOfSetup.X = withRoutingCall
}

// appendToExistingRouting appends a route statement to an existing WithRouting call.
func appendToExistingRouting(withRoutingCall *dst.CallExpr, routeStmt *dst.ExprStmt) {
	if len(withRoutingCall.Args) == 0 {
		return
	}

	funcLit, ok := withRoutingCall.Args[0].(*dst.FuncLit)
	if !ok {
		return
	}

	// Check if route already exists
	for _, stmt := range funcLit.Body.List {
		if match.EqualNode(routeStmt).MatchNode(stmt) {
			return
		}
	}

	// Ensure proper formatting for the first statement
	if len(funcLit.Body.List) > 0 {
		// Ensure existing statements have proper newlines
		if funcLit.Body.List[0].Decorations().Before != dst.NewLine {
			funcLit.Body.List[0].Decorations().Before = dst.NewLine
		}
	}

	// Add proper decorations to the new route statement
	routeStmt.Decs.Before = dst.NewLine
	routeStmt.Decs.After = dst.NewLine

	funcLit.Body.List = append(funcLit.Body.List, routeStmt)
}

// removeRouteFromFuncBody removes a route statement from the function body.
func removeRouteFromFuncBody(body *dst.BlockStmt, routeStmt *dst.ExprStmt) {
	body.List = slices.DeleteFunc(body.List, func(stmt dst.Stmt) bool {
		return match.EqualNode(routeStmt).MatchNode(stmt)
	})
}
