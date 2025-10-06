# Scheduling Report System - Go Configuration API

> REST API for managing report configurations, schedules, deliveries, and recipients (Phases 1-6)

## 🎯 Purpose

This Go API serves as the **Configuration Management Layer** of the Scheduling Report System. It provides a complete REST API for:

- Managing data sources (MySQL, PostgreSQL, Oracle, SQL Server, MongoDB, BigQuery, Snowflake)
- Creating and versioning report configurations
- Scheduling reports with cron expressions
- Configuring delivery methods (Email, SFTP, Webhook, S3, File Share)
- Managing delivery recipients
- Tracking execution history
- Auditing configuration changes

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0+
- Access to target data sources

### Installation

```bash
# Install dependencies
go mod download

# Copy environment configuration
cp .env.example .env

# Edit .env with your database credentials
nano .env
```

### Configuration

Edit `.env`:

```env
# Server Configuration
APP_PORT=3000

# Database Configuration
DB_HOST=192.168.131.209
DB_PORT=3306
DB_USER=user_arif
DB_PASSWORD=your_password
DB_NAME=lab
```

### Running the API

```bash
# Build and run
go build -o scheduling-report && ./scheduling-report

# Or run directly
go run main.go
```

The API will start on `http://localhost:3000`

## 📋 API Endpoints

### Phase 1: Data Sources (5 endpoints)
- `GET    /api/datasources` - List all data sources
- `GET    /api/datasources/:id` - Get data source by ID
- `POST   /api/datasources` - Create new data source
- `PUT    /api/datasources/:id` - Update data source
- `DELETE /api/datasources/:id` - Soft delete data source

### Phase 2: Report Configs (5 endpoints)
- `GET    /api/report-configs` - List all report configs
- `GET    /api/report-configs/:id` - Get report config by ID
- `POST   /api/report-configs` - Create report config
- `PUT    /api/report-configs/:id` - Update report config (increments version)
- `DELETE /api/report-configs/:id` - Soft delete report config

### Phase 3: Schedules (6 endpoints)
- `GET    /api/schedules` - List all schedules
- `GET    /api/schedules/:id` - Get schedule by ID
- `GET    /api/schedules/config/:config_id` - Get schedules by config
- `POST   /api/schedules` - Create schedule
- `PUT    /api/schedules/:id` - Update schedule
- `DELETE /api/schedules/:id` - Soft delete schedule

### Phase 4: Deliveries & Recipients (12 endpoints)

**Deliveries:**
- `GET    /api/deliveries` - List all deliveries
- `GET    /api/deliveries/:id` - Get delivery by ID
- `GET    /api/deliveries/config/:config_id` - Get deliveries by config
- `POST   /api/deliveries` - Create delivery
- `PUT    /api/deliveries/:id` - Update delivery
- `DELETE /api/deliveries/:id` - Soft delete delivery

**Recipients:**
- `GET    /api/recipients` - List all recipients
- `GET    /api/recipients/:id` - Get recipient by ID
- `GET    /api/recipients/delivery/:delivery_id` - Get recipients by delivery
- `POST   /api/recipients` - Create recipient
- `PUT    /api/recipients/:id` - Update recipient
- `DELETE /api/recipients/:id` - Hard delete recipient

### Phase 5: Executions & Logs (7 endpoints)

**Executions:**
- `GET    /api/executions` - List all executions
- `GET    /api/executions/:id` - Get execution by ID
- `GET    /api/executions/config/:config_id` - Get executions by config

**Delivery Logs:**
- `GET    /api/delivery-logs` - List all delivery logs
- `GET    /api/delivery-logs/:id` - Get delivery log by ID
- `GET    /api/delivery-logs/execution/:execution_id` - Get logs by execution
- `GET    /api/delivery-logs/delivery/:delivery_id` - Get logs by delivery

### Phase 6: Audits (4 endpoints)
- `GET    /api/audits` - List all audits
- `GET    /api/audits/:id` - Get audit by ID
- `GET    /api/audits/config/:config_id` - Get audits by config
- `GET    /api/audits/user/:user_id` - Get audits by user

**Total: 39 endpoints**

## 🏗️ Project Structure

```
golang-api/
├── main.go                 # Application entry point
├── config/
│   └── database.go         # Database connection setup
├── models/                 # Database models (GORM)
│   ├── datasource.go
│   ├── report_config.go
│   ├── report_schedule.go
│   ├── report_delivery.go
│   ├── report_delivery_recipient.go
│   ├── report_execution.go
│   ├── report_delivery_log.go
│   └── report_config_audit.go
├── repositories/           # Data access layer
│   ├── datasource_repository.go
│   ├── report_config_repository.go
│   ├── report_schedule_repository.go
│   ├── report_delivery_repository.go
│   ├── report_delivery_recipient_repository.go
│   ├── report_execution_repository.go
│   ├── report_delivery_log_repository.go
│   └── report_config_audit_repository.go
├── services/               # Business logic layer
│   ├── datasource_service.go
│   ├── report_config_service.go
│   ├── report_schedule_service.go
│   ├── report_delivery_service.go
│   ├── report_delivery_recipient_service.go
│   ├── report_execution_service.go
│   ├── report_delivery_log_service.go
│   └── report_config_audit_service.go
├── controllers/            # HTTP handlers
│   ├── datasource_controller.go
│   ├── report_config_controller.go
│   ├── report_schedule_controller.go
│   ├── report_delivery_controller.go
│   ├── report_delivery_recipient_controller.go
│   ├── report_execution_controller.go
│   ├── report_delivery_log_controller.go
│   └── report_config_audit_controller.go
├── routes/
│   └── route.go            # Route definitions
├── middlewares/
│   └── logger.go           # Request logging
└── utils/
    └── response.go         # Standard response format
```

## 🔧 Key Technologies

- **Framework:** Fiber v2 (high-performance web framework)
- **ORM:** GORM (database abstraction)
- **Database:** MySQL 8.0+
- **Cron Parser:** robfig/cron/v3 (schedule validation)
- **Environment:** godotenv (configuration management)

## 💾 Database Schema

The API manages 8 MySQL tables:

1. **report_datasources** - Data source connections
2. **report_configs** - Report definitions with versioning
3. **report_schedules** - Cron-based scheduling
4. **report_deliveries** - Delivery method configurations
5. **report_delivery_recipients** - Recipient lists
6. **report_executions** - Execution history (UUID-based)
7. **report_delivery_logs** - Delivery attempt logs
8. **report_config_audits** - Automatic audit trail

See `/docs/REPORT-HUB.md` for complete schema documentation.

## 📝 Response Format

All API responses follow this standard format:

### Success Response
```json
{
  "code": "20000000",
  "status": "success",
  "message": "Operation successful",
  "data": { /* response data */ }
}
```

### Error Response
```json
{
  "code": "40000001",
  "status": "error",
  "message": "Error description",
  "error": "Detailed error message"
}
```

### Response Code Format
`{HTTP_STATUS}{SERVICE_CODE}{ERROR_CODE}`

- **HTTP_STATUS:** 200, 400, 404, 500, etc.
- **SERVICE_CODE:** 00 (General), 01 (Datasource), 02 (Config), 03 (Schedule), 04 (Delivery), 05 (Recipient), 06 (Execution), 07 (Log), 08 (Audit)
- **ERROR_CODE:** 00 (Success), 01 (Validation), 02 (Not Found), 03 (Duplicate), 04 (FK Violation), etc.

## 🔍 Key Features

### Automatic Audit Trail
All create/update/delete operations on report configs automatically create audit records:

```go
// Audit is automatically created in service layer
audit := models.ReportConfigAudit{
    ConfigID:   config.ID,
    ActionType: "UPDATE",
    OldValue:   oldConfigJSON,
    NewValue:   newConfigJSON,
    ChangedBy:  userId,
}
```

### Foreign Key Validation
Services validate all foreign key references before operations:

```go
// Check if datasource exists before creating report config
if err := s.datasourceRepo.CheckDatasourceExists(configReq.DatasourceID); err != nil {
    return nil, "40002004", "Invalid datasource_id"
}
```

### Soft Delete Pattern
Most entities use soft delete with `is_active` flag:

```go
// Soft delete - sets is_active = false
config.IsActive = false
config.DeletedAt = &now
config.DeletedBy = &userId
```

### Version Control
Report configs auto-increment version on updates:

```go
existingConfig.Version++  // Auto-incremented
```

### Cron Validation
Schedule cron expressions are validated:

```go
if _, err := cron.ParseStandard(scheduleReq.CronExpression); err != nil {
    return nil, "40003001", "Invalid cron expression format"
}
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./services/...
```

## 🚢 Deployment

### Build for Production

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scheduling-report .

# Run production binary
./scheduling-report
```

### Docker Deployment

```bash
# Build Docker image
docker build -t scheduling-report-api .

# Run container
docker run -p 3000:3000 --env-file .env scheduling-report-api
```

## 📚 Documentation

- **API Specification:** `/docs/api-docs.yaml` (OpenAPI 3.0)
- **Database Schema:** `/docs/REPORT-HUB.md`
- **Project Status:** `/docs/PROJECT_STATUS.md`
- **Python Integration:** `/docs/PYTHON_PHASES_7_8.md`

## 🔗 Integration with Python API

This Go API works alongside the Python Execution Engine:

- **Go API:** Configuration management (this service)
- **Python API:** Report execution and scheduling (Phases 7-8)
- **Shared Database:** Both services access the same MySQL database
- **Communication:** Python reads configs, writes executions/logs

See `/docs/PYTHON_PHASES_7_8.md` for Python implementation details.

## 🐛 Common Issues

### Database Connection Fails
```bash
# Check MySQL is accessible
mysql -h192.168.131.209 -uuser_arif -p lab

# Verify .env credentials match
```

### Port Already in Use
```bash
# Change APP_PORT in .env
APP_PORT=3001

# Or kill existing process
pkill -f scheduling-report
```

### GORM Auto-Migration Issues
```bash
# Run manual migrations from /docs/migrations/
mysql -h192.168.131.209 -uuser_arif -p lab < ../docs/migrations/001_initial_schema.sql
```

## 📊 Current Statistics

| Metric | Count |
|--------|-------|
| Total Endpoints | 39 |
| Database Tables | 8 |
| Supported Data Sources | 7 |
| Delivery Methods | 5 |
| Service Files | 8 |
| Controller Files | 8 |
| Repository Files | 8 |
| Model Files | 8 |

## 📄 License

This project is part of the Scheduling Report System.

## 🤝 Related Services

- **Python Execution Engine:** `/python-api/` (Phases 7-8)
- **Documentation:** `/docs/`

---

**Status:** Phase 1-6 Complete ✅ | **Next:** Python Integration (Phases 7-8)
