package configs

import "os"

// Функция для получения значения из переменной окружения или использования значения по умолчанию
func GetEnvOrDefault(flagValue string, envKey, defaultValue string) string {
	// Если флаг не пустой, используем его значение
	if flagValue != "" {
		return flagValue
	}

	// Если переменная окружения существует, используем ее значение
	if envValue, exists := os.LookupEnv(envKey); exists {
		return envValue
	}

	// Если переменная окружения не существует, возвращаем значение по умолчанию
	return defaultValue
}
