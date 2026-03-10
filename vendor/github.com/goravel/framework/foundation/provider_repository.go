package foundation

import (
	"fmt"
	"sort"
	"strings"

	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/errors"
)

type ProviderState struct {
	provider   foundation.ServiceProvider
	registered bool
	booted     bool
}

var _ foundation.ProviderRepository = (*ProviderRepository)(nil)

type ProviderRepository struct {
	states      map[string]*ProviderState
	providers   []foundation.ServiceProvider
	sorted      []foundation.ServiceProvider
	sortedValid bool
	loaded      bool
}

func NewProviderRepository() *ProviderRepository {
	return &ProviderRepository{
		states:    make(map[string]*ProviderState),
		providers: make([]foundation.ServiceProvider, 0),
	}
}

func (r *ProviderRepository) Add(providers []foundation.ServiceProvider) {
	if len(providers) == 0 {
		return
	}

	for _, provider := range providers {
		key := r.getProviderName(provider)

		if _, exists := r.states[key]; exists {
			continue
		}

		state := &ProviderState{provider: provider}
		r.states[key] = state
		r.providers = append(r.providers, provider)
	}

	r.sortedValid = false
	r.sorted = nil
}

func (r *ProviderRepository) Boot(app foundation.Application) {
	providers := r.getProviders()

	for _, provider := range providers {
		key := r.getProviderName(provider)
		state, exists := r.states[key]

		if !exists || !state.registered || state.booted {
			continue
		}

		state.provider.Boot(app)
		state.booted = true
	}
}

func (r *ProviderRepository) GetBooted() []foundation.ServiceProvider {
	booted := make([]foundation.ServiceProvider, 0, len(r.states))
	for _, state := range r.states {
		if state.booted {
			booted = append(booted, state.provider)
		}
	}
	return booted
}

func (r *ProviderRepository) LoadFromConfig(config config.Config) []foundation.ServiceProvider {
	if r.loaded {
		return r.providers
	}

	if config == nil {
		return r.providers
	}

	raw := config.Get("app.providers")
	providers, ok := raw.([]foundation.ServiceProvider)
	if !ok {
		return r.providers
	}

	r.Add(providers)
	r.loaded = true
	return r.providers
}

func (r *ProviderRepository) Register(app foundation.Application) []foundation.ServiceProvider {
	providers := r.getProviders()
	processed := make([]foundation.ServiceProvider, 0, len(providers))

	for _, provider := range providers {
		key := r.getProviderName(provider)
		state, exists := r.states[key]
		if !exists {
			state = &ProviderState{provider: provider}
			r.states[key] = state
		}

		if state.registered {
			processed = append(processed, provider)
			continue
		}

		provider.Register(app)
		state.registered = true
		processed = append(processed, provider)
	}

	return processed
}

func (r *ProviderRepository) Reset() {
	r.providers = make([]foundation.ServiceProvider, 0)
	r.states = make(map[string]*ProviderState)
	r.sorted = make([]foundation.ServiceProvider, 0)
	r.sortedValid = false
	r.loaded = false
}

func (r *ProviderRepository) getProviders() []foundation.ServiceProvider {
	if r.sortedValid {
		return r.sorted
	}

	r.sorted = r.sort(r.providers)
	r.sortedValid = true
	return r.sorted
}

func (r *ProviderRepository) getRelationship(provider foundation.ServiceProvider) binding.Relationship {
	if p, ok := provider.(foundation.ServiceProviderWithRelations); ok {
		return p.Relationship()
	}
	return binding.Relationship{}
}

func (r *ProviderRepository) getProviderName(provider foundation.ServiceProvider) string {
	return fmt.Sprintf("%T", provider)
}

// sort performs a topological sort on a list of providers to ensure
// providers with dependencies are registered and booted *after*
// the providers they depend on.
func (r *ProviderRepository) sort(providers []foundation.ServiceProvider) []foundation.ServiceProvider {
	if len(providers) == 0 {
		return providers
	}

	// These maps build a directed graph of dependencies.
	bindingToProvider := make(map[string]foundation.ServiceProvider)
	providerToVirtualBinding := make(map[foundation.ServiceProvider]string)
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	virtualBindingCounter := 0

	// --- Build Graph Nodes ---
	// Identify all "nodes" (bindings) in the graph.
	for _, provider := range providers {
		relationship := r.getRelationship(provider)
		bindings := relationship.Bindings
		dependencies := relationship.Dependencies
		provideFor := relationship.ProvideFor

		if len(bindings) > 0 {
			for _, b := range bindings {
				bindingToProvider[b] = provider
				inDegree[b] = 0
			}
		} else if len(dependencies) > 0 || len(provideFor) > 0 {
			// This provider has no bindings but has relationships.
			// We create a "virtual" node to represent it in the graph
			// so its dependencies can be sorted.
			virtualBinding := fmt.Sprintf("__virtual_%d", virtualBindingCounter)
			virtualBindingCounter++
			bindingToProvider[virtualBinding] = provider
			providerToVirtualBinding[provider] = virtualBinding
			inDegree[virtualBinding] = 0
		}
	}

	// --- Build Graph Edges ---
	// Connect the nodes based on 'Dependencies' and 'ProvideFor'.
	for _, provider := range providers {
		relationship := r.getRelationship(provider)
		bindings := relationship.Bindings
		dependencies := relationship.Dependencies
		provideFor := relationship.ProvideFor

		var providerBindings []string
		if len(bindings) > 0 {
			providerBindings = bindings
		} else if virtualBinding, exists := providerToVirtualBinding[provider]; exists {
			providerBindings = []string{virtualBinding}
		}

		if len(providerBindings) == 0 {
			// Provider is independent and not part of the sort.
			continue
		}

		for _, b := range providerBindings {
			// Edge: dep -> binding
			for _, dep := range dependencies {
				if _, exists := bindingToProvider[dep]; exists {
					graph[dep] = append(graph[dep], b)
					inDegree[b]++
				}
			}

			// Edge: binding -> provideForBinding
			for _, provideForBinding := range provideFor {
				if _, exists := bindingToProvider[provideForBinding]; exists {
					graph[b] = append(graph[b], provideForBinding)
					inDegree[provideForBinding]++
				}
			}
		}
	}

	// --- Topological Sort (Kahn's Algorithm) ---
	queue := make([]string, 0, len(inDegree))
	for b, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, b)
		}
	}

	result := make([]string, 0, len(inDegree))
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// --- Cycle Detection & Result Reconstruction ---
	if len(result) != len(inDegree) {
		// A cycle exists, the dependency order is impossible.
		cycle := r.detectCycle(graph, bindingToProvider)
		if len(cycle) > 0 {
			panic(errors.ServiceProviderCycle.Args(strings.Join(cycle, " -> ")))
		}
		panic(errors.ServiceProviderCycle.Args("unknown cycle detected"))
	}

	sortedProviders := make([]foundation.ServiceProvider, 0, len(providers))
	used := make(map[foundation.ServiceProvider]bool)

	for _, b := range result {
		provider := bindingToProvider[b]
		// Use a map to prevent adding the same provider multiple
		// times if it had more than one binding (e.g., "log" and "logger").
		if !used[provider] {
			sortedProviders = append(sortedProviders, provider)
			used[provider] = true
		}
	}

	// Add any remaining providers that were not part of the graph.
	for _, provider := range providers {
		if !used[provider] {
			sortedProviders = append(sortedProviders, provider)
		}
	}

	return sortedProviders
}

// detectCycle uses a Depth-First Search (DFS) to find and report a
// cycle in the provider dependency graph.
func (r *ProviderRepository) detectCycle(graph map[string][]string, bindingToProvider map[string]foundation.ServiceProvider) []string {
	// visited: Nodes already processed and known to be safe.
	// recStack: Nodes currently in our DFS recursion stack.
	// If we hit a node already in recStack, we've found a cycle.
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := make([]string, 0)
	cycle := make([]string, 0)

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				// Cycle detected! Reconstruct the path for the error message.
				cycleStart := -1
				for i, p := range path {
					if p == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					cycle = append(cycle, path[cycleStart:]...)
					// Append the start node to the end to show the full loop.
					cycle = append(cycle, neighbor)
				}
				return true
			}
		}

		// Backtrack
		recStack[node] = false
		path = path[:len(path)-1]
		return false
	}

	// Sort nodes to make cycle detection deterministic.
	var nodes []string
	for node := range graph {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	for _, node := range nodes {
		if !visited[node] {
			if dfs(node) {
				break
			}
		}
	}

	if len(cycle) == 0 {
		return nil
	}

	// Convert the list of *bindings* into user-friendly *provider names*.
	var cycleProviders []string
	providerSet := make(map[string]struct{})

	for _, b := range cycle {
		if provider, exists := bindingToProvider[b]; exists {
			providerName := r.getProviderName(provider)
			cycleProviders = append(cycleProviders, providerName)
			providerSet[providerName] = struct{}{}
		}
	}

	// Handle a specific edge case for self-loops.
	if len(cycleProviders) == 2 && cycleProviders[0] == cycleProviders[1] {
		if len(providerSet) == 1 && len(cycle) > 2 {
			return cycleProviders[0:1]
		}
	}

	return cycleProviders
}
