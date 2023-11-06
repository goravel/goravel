package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("hashing", map[string]any{
		// Hashing Driver
		//
		// This option controls the default diver that gets used
		// by the framework hash facade.
		// Default driver is "bcrypt".
		//
		// Supported Drivers: "bcrypt", "argon2id"
		"driver": "bcrypt",

		// Bcrypt Hashing Options
		// rounds: The cost factor that should be used to compute the bcrypt hash.
		// The cost factor controls how much time is needed to compute a single bcrypt hash.
		// The higher the cost factor, the more hashing rounds are done. Increasing the cost
		// factor by 1 doubles the necessary time. After a certain point, the returns on
		// hashing time versus attacker time are diminishing, so choose your cost factor wisely.
		"bcrypt": map[string]any{
			"rounds": 12,
		},

		// Argon2id Hashing Options
		// memory: A memory cost, which defines the memory usage, given in kibibytes.
		// time: A time cost, which defines the amount of computation
		// realized and therefore the execution time, given in number of iterations.
		// threads: A parallelism degree, which defines the number of parallel threads.
		"argon2id": map[string]any{
			"memory":  65536,
			"time":    4,
			"threads": 1,
		},
	})
}
