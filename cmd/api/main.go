package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	// Importe o pacote docs sem o "_" para poder acessar a variável SwaggerInfo
	"github.com/kisalto/Feel-The-Night-Swagger/docs"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/app"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/database"
)

// @title           Feel The Night API
// @version         1.0
// @description     API para gerenciamento de guias, eventos e personagens.
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	if err := godotenv.Load(); err != nil {
		zap.L().Warn("Arquivo .env não encontrado, lendo variáveis de ambiente do sistema.")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		zap.L().Fatal("JWT_SECRET não configurado")
	}

	if len(secret) < 32 {
		zap.L().Fatal("JWT_SECRET precisa ter pelo menos 32 caracteres")
	}

	// 1. Tratamento da Porta
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	// 2. Configuração Dinâmica do Swagger via docs.SwaggerInfo
	renderHost := os.Getenv("RENDER_EXTERNAL_HOSTNAME")

	if renderHost != "" {
		// Em produção no Render
		docs.SwaggerInfo.Host = renderHost
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		// Em desenvolvimento local
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	// 3. Conexão com o Banco de Dados
	zap.L().Info("Iniciando conexão com o banco de dados...")
	db, err := database.Connect()
	if err != nil {
		zap.L().Fatal("Erro crítico ao inicializar o banco", zap.Error(err))
	}

	// 4. Carrega Rotas
	routes := app.NewRouter(db)

	// 5. Inicia o Servidor
	serverPort := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:              serverPort,
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		zap.L().Info("Servidor rodando na porta", zap.String("port", serverPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("Erro ao iniciar o servidor", zap.Error(err))
		}
	}()

	<-ctx.Done()
	zap.L().Info("Sinal recebido, desligando o servidor...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("Shutdown forçado: requisições não terminaram a tempo", zap.Error(err))
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}

	zap.L().Info("Servidor desligado")

}
