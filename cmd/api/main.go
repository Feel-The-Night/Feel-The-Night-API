package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	// Importe o pacote docs sem o "_" para poder acessar a variável SwaggerInfo
	"github.com/kisalto/Feel-The-Night-Swagger/docs"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/database"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/handler"
)

// @title           Feel The Night API
// @version         1.0
// @description     API para gerenciamento de guias, eventos e personagens.
// @BasePath        /
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
	if err := database.Connect(); err != nil {
		zap.L().Fatal("Erro crítico ao inicializar o banco", zap.Error(err))
	}

	// 4. Carrega Rotas
	routes := handler.Setup()

	// 5. Inicia o Servidor
	serverPort := fmt.Sprintf(":%s", port)
	zap.L().Info("Servidor rodando na porta", zap.String("port", serverPort))

	if err := http.ListenAndServe(serverPort, routes); err != nil {
		zap.L().Fatal("Erro ao iniciar o servidor", zap.Error(err))
	}
}
