package logger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Init configura el logger de sistema global
func Init(env string) {
	var handler slog.Handler

	// Configuramos Lumberjack para la rotación de archivos
	fileLogger := &lumberjack.Logger{
		Filename:   "logs/app.log", // Ruta donde se guardarán
		MaxSize:    10,             // Tamaño máximo en MB antes de rotar
		MaxBackups: 5,              // Cuántos archivos viejos conservar
		MaxAge:     28,             // Días máximos a conservar antes de borrar
		Compress:   true,           // Comprimir archivos rotados (.gz)
	}

	if env == "production" {
		// Output para Producción: Archivo rotativo + Consola (JSON)
		multiWriter := io.MultiWriter(os.Stdout, fileLogger)
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Output para Desarrollo local: Archivo rotativo (JSON) + Consola local (Texto legible)

		// Opcional para dev puro: si solo quieres consola, podrías dejar solo os.Stdout
		multiWriter := io.MultiWriter(os.Stdout, fileLogger)
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)

	// Establecemos este logger estructurado como el default en todo Go
	slog.SetDefault(logger)
}
