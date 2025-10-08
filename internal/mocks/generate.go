package mocks

//go:generate mockgen -source=../service/mcp_service.go -destination=mcp_service_mock.go -package=mocks
//go:generate mockgen -source=../repository/user_interfaces.go -destination=user_repository_mock.go -package=mocks
//go:generate mockgen -source=../mcp/tool.go -destination=mcp_tool_mock.go -package=mocks
//go:generate mockgen -source=../provider/types.go -destination=provider_mock.go -package=mocks
//go:generate mockgen -source=../googleai/client.go -destination=googleai_client_mock.go -package=mocks
//go:generate mockgen -source=../openai/client.go -destination=openai_client_mock.go -package=mocks