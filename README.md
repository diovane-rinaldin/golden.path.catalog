# Golden Path Portal

## Descrição
Um portal inspirado no CNCF Landscape para gerenciar e visualizar as ferramentas e tecnologias mantidas pela empresa. O portal permite o cadastro de tecnologias, componentes e serviços, organizados de forma hierárquica e com suporte a imagens.

## Arquitetura

### Backend (Golang)
* Framework Gin para APIs REST
* Autenticação JWT com chaves armazenadas no AWS KMS
* Armazenamento de dados no DynamoDB
* Upload de imagens para S3
* Logs com OpenTelemetry
* Estrutura MVC
* Containerização com Docker
* Deploy com Kubernetes

### Frontend (React)
* Interface responsiva com Tailwind CSS
* Gerenciamento de estado com Context API
* Upload de imagens
* Navegação com React Router
* Containerização com Docker
* Deploy com Kubernetes

### Infraestrutura (Terraform)
* Provisionamento do DynamoDB
* Configuração do bucket S3
* Gerenciamento de chaves KMS
* IAM Roles e Policies

## Estrutura de Dados

### Technology
```json
{
    "id": "string",
    "name": "string",
    "description": "string",
    "image_url": "string"
}
```

### Component
```json
{
    "id": "string",
    "technology_id": "string",
    "name": "string",
    "description": "string",
    "image_url": "string"
}
```

### Service
```json
{
    "id": "string",
    "component_id": "string",
    "name": "string",
    "description": "string",
    "image_url": "string",
    "cloud_provider": "string",
    "service_cloud_url": "string"
}
```

## Pré-requisitos

* Go 1.20+
* Node.js 18+
* Docker
* Kubernetes cluster
* AWS CLI configurado
* Terraform
* Make (opcional)

## Configuração e Instalação

### Infraestrutura

1. Configure suas credenciais AWS:
```bash
aws configure
```

2. Instale e configure a infraestrutura:
```bash
cd infrastructure
terraform init
terraform plan
terraform apply
```

Caso precise instalar o Terraform
```bash
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo
sudo yum -y install terraform
```

**Problemas conhecidos**

Eventualmente ao criar o bucket pode haver um erro a respeito das políticas públicas

```bash
│ Error: putting S3 Bucket (golden-path-images) Policy: operation error S3: PutBucketPolicy, https response error StatusCode: 403, RequestID: G1XC8S177HW5ETQ1, HostID: 9cKOXV0Pp6D+g62z5EhsosCoqAAkQkNz4Yrh1WcN7RlGT8wCXTQKRV3I+XxdSkk91WxODPvXY2I=, api error AccessDenied: User: arn:aws:iam::948052180351:root is not authorized to perform: s3:PutBucketPolicy on resource: "arn:aws:s3:::golden-path-images" because public policies are blocked by the BlockPublicPolicy block public access setting.
│ 
│   with aws_s3_bucket_policy.public_read,
│   on main.tf line 108, in resource "aws_s3_bucket_policy" "public_read":
│  108: resource "aws_s3_bucket_policy" "public_read" {
│ 
╵
```
Para corrigir, crie o bucket manualmente e execute novamente o `terraform apply`
```bash
aws s3api create-bucket --bucket golden-path-images --region us-east-1
aws s3api get-public-access-block --bucket golden-path-images
aws s3api delete-public-access-block --bucket golden-path-images
terraform apply
```

3. Guarde os outputs gerados para configurar as aplicações:
```bash
terraform output
```

### Backend

1. Configure o ambiente:
```bash
cd backend
cp .env.example .env
# Edite o arquivo .env com os valores obtidos do Terraform
```

2. Instale as dependências:
```bash
go mod download
```

3. Build da imagem Docker:
```bash
docker build -t golden-path-backend .
```

4. Deploy no Kubernetes:
```bash
kubectl apply -f backend-deployment.yaml
```

### Frontend

1. Configure o ambiente:
```bash
cd frontend
cp .env.example .env
# Edite o arquivo .env com os valores apropriados
```

2. Instale as dependências:
```bash
npm install
```

3. Build da imagem Docker:
```bash
docker build -t golden-path-frontend .
```

4. Deploy no Kubernetes:
```bash
kubectl apply -f frontend-deployment.yaml
```

## Executando Localmente

### Infraestrutura
```bash
cd infrastructure
terraform init
terraform plan    # visualizar mudanças
terraform apply   # aplicar mudanças
terraform destroy # destruir infraestrutura
```

### Backend
```bash
cd backend
# Execute diretamente
go run main.go

# Ou via Docker
docker run -p 8080:8080 --env-file .env golden-path-backend
```

O servidor estará disponível em `http://localhost:8080`

### Frontend
```bash
cd frontend
# Execute diretamente
npm start

# Ou via Docker
docker run -p 8081:8081 golden-path-frontend
```

A aplicação estará disponível em `http://localhost:8081`

Erros conhecidos
* **Module not found: Error: Can't resolve 'axios':** npm install axios --save --force --legacy-peer-deps
* **Module not found: Error: Can't resolve 'react-router-dom':** npm install react-router-dom --save --force --legacy-peer-deps

## APIs Disponíveis

### Autenticação
* `POST /auth` - Autenticação e geração de token

### Technology
* `POST /api/technology` - Criar nova tecnologia
* `GET /api/technology` - Listar todas as tecnologias
* `GET /api/technology/:name` - Buscar tecnologia por nome
* `PUT /api/technology/:id` - Atualizar tecnologia

### Component
* `POST /api/component` - Criar novo componente
* `GET /api/component` - Listar todos os componentes
* `GET /api/component/:name` - Buscar componente por nome
* `GET /api/component/technology/:id` - Listar componentes por tecnologia
* `PUT /api/component/:id` - Atualizar componente

### Service
* `POST /api/service` - Criar novo serviço
* `GET /api/service` - Listar todos os serviços
* `GET /api/service/:name` - Buscar serviço por nome
* `GET /api/service/component/:id` - Listar serviços por componente
* `PUT /api/service/:id` - Atualizar serviço

### Upload
* `POST /api/upload` - Upload de imagem para S3

## Variáveis de Ambiente

### Backend (.env)
```
AWS_REGION=us-east-1
DYNAMODB_ENDPOINT=https://dynamodb.us-east-1.amazonaws.com
S3_BUCKET_NAME=golden_path_images
S3_BUCKET_URL=https://golden_path_images.s3.amazonaws.com
KMS_KEY_ID=alias/golden_path_api_auth
JWT_SECRET=your-secret-key-here
```

### Frontend (.env)
```
REACT_APP_API_URL=http://localhost:8080
REACT_APP_AWS_REGION=us-east-1
REACT_APP_KMS_KEY_ID=alias/golden_path_api_auth
```

## Monitoramento e Logs

* Os logs do backend são gerados usando OpenTelemetry
* Métricas do Kubernetes disponíveis via Prometheus
* Logs dos containers disponíveis via kubectl logs

## Contribuindo

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/amazing-feature`)
3. Commit suas mudanças (`git commit -m 'Add some amazing feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE.md](LICENSE.md) para detalhes.