# Performance Testing: DB vs ORM

This directory contains performance benchmarks to compare direct database operations (using `facades.DB()`) versus ORM operations (using `facades.Orm()`) in the Goravel framework.

## Background

The issue ([#754](https://github.com/goravel/goravel/issues/754)) reports that DB operations should be faster than ORM operations based on FrameworkBenchmarks results, but this doesn't appear to be consistently true. These tests provide a way to measure and compare performance locally.

## Running the Tests

### Prerequisites

1. PostgreSQL running locally on port 5432
2. Database named "goravel" 
3. Database user "root" with access to the database

### Basic Setup

1. Start PostgreSQL:
   ```bash
   # Example using Docker
   docker run --name postgres -e POSTGRES_DB=goravel -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
   ```

2. Update `.env` file with database credentials:
   ```
   DB_CONNECTION=postgres
   DB_HOST=127.0.0.1
   DB_PORT=5432
   DB_DATABASE=goravel
   DB_USERNAME=root
   DB_PASSWORD=password
   ```

3. Run migrations to create the users table:
   ```bash
   go run . artisan migrate
   ```

### Running Tests

#### Quick Performance Comparison
```bash
# Run basic test with performance comparison
go test -v ./tests/performance/
```

#### Detailed Benchmarks
```bash
# Run all benchmarks
go test -bench=. ./tests/performance/

# Run specific benchmarks
go test -bench=BenchmarkDBSelect ./tests/performance/
go test -bench=BenchmarkORMSelect ./tests/performance/

# Run with memory profiling
go test -bench=. -benchmem ./tests/performance/

# Run with CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./tests/performance/
```

#### Benchmark Options
```bash
# Run each benchmark for specific duration
go test -bench=. -benchtime=10s ./tests/performance/

# Run with specific number of iterations
go test -bench=. -count=5 ./tests/performance/

# Output results to file
go test -bench=. ./tests/performance/ > benchmark_results.txt
```

## What's Being Tested

### SELECT Operations
- **DB**: `facades.DB().Table("users").Limit(10).Get(&results)`
- **ORM**: `facades.Orm().Query().Limit(10).Find(&users)`

### INSERT Operations  
- **DB**: `facades.DB().Table("users").Insert(data)`
- **ORM**: `facades.Orm().Query().Create(&user)`

### UPDATE Operations
- **DB**: `facades.DB().Table("users").Where("id", id).Update(data)`
- **ORM**: `facades.Orm().Query().Where("id", id).Update("name", value)`

## Performance Metrics

The benchmarks measure:
- **Execution Time**: nanoseconds per operation
- **Memory Allocations**: bytes allocated per operation  
- **Throughput**: operations per second

## Expected Results

Generally, direct DB operations should be faster than ORM operations because:
1. **Less Abstraction**: DB operations have fewer layers
2. **No Model Mapping**: No need to map results to struct instances
3. **Simpler Query Building**: Direct query building vs ORM query builder

However, the performance difference depends on:
- Query complexity
- Data size
- Database configuration
- Network latency

## Interpreting Results

Example benchmark output:
```
BenchmarkDBSelect-8      1000    1234567 ns/op    1024 B/op    10 allocs/op
BenchmarkORMSelect-8      500    2345678 ns/op    2048 B/op    20 allocs/op
```

This means:
- DB SELECT: 1000 iterations, ~1.23ms per operation, 1KB allocated per operation
- ORM SELECT: 500 iterations, ~2.35ms per operation, 2KB allocated per operation
- ORM is ~1.9x slower than DB for SELECT operations

## Troubleshooting

### "Database not available" 
- Ensure PostgreSQL is running
- Verify database credentials in `.env`
- Check if the database and tables exist

### Connection Refused
```bash
# Check if PostgreSQL is running
pg_ctl status

# Start PostgreSQL if needed
pg_ctl start
```

### Table Doesn't Exist
```bash
# Run migrations
go run . artisan migrate

# Or manually create users table
go run . artisan migrate:fresh
```

## Contributing

When adding new benchmarks:
1. Follow the naming convention: `BenchmarkDB*` and `BenchmarkORM*`
2. Include proper error handling for missing database/tables
3. Use consistent measurement patterns
4. Add memory profiling with `runtime.MemStats`
5. Document what's being tested

## Advanced Usage

### Comparing Different Query Patterns
You can extend the benchmarks to test:
- Complex WHERE clauses
- JOIN operations  
- Bulk inserts
- Transaction performance
- Connection pooling effects

### Custom Database Configuration
Test with different database configurations:
- Connection pool sizes
- Query timeouts  
- SSL settings
- Different PostgreSQL versions