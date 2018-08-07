### Maestro

Maestro é o componente de plataforma responsável por orquestrar a parte de persistência no domain, levando em consideração a necessidade de reprocessamento. Além disso, ele também faz a orquerstração de todo o processo de reprocessamento

### Build

O build da aplicação é feito através de um arquivo Makefile, para buildar a aplicação execute o seguinte comando:

```sh
$ make
```

Após executar o make será criada uma pasta dist e o executável da aplicação maestro.

### Deploy

O processo de deploy do maestro na plataforma é feito através do installer, os componentes em Go são compilados e comitados dentro do installer então para atualizar a versão do maestro para atualizar a versão do maestro na plataforma utilize o seguinte comando:

```sh
$ mv dist/maestro ~/installed_plataforma/Plataforma-Installer/Dockerfiles
$ plataforma --upgrade maestro
```

### API

Retornar a lista de reprocessamentos pendentes por sitema
```http
GET /v1.0.0/reprocessing/<solution_id>/pending HTTP/1.1
Host: localhost:6971
```

Aprovar um reprocessamento
```http
POST /v1.0.0/reprocessing/<reprocessing_id>/approve HTTP/1.1
Host: localhost:6971
Content-Type: application/json

{
  "user":"<usuário aprovador>"
}
```

Ignorar um reprocessamento

```http
POST /v1.0.0/reprocessing/<reprocessing_id>/skip HTTP/1.1
Host: localhost:6971
Content-Type: application/json

{
  "user":"<usuário que ignorou>"
}
```

Verificar se a plataforma está habilitada para receber eventos

```http
GET /v1.0.0/gateway/<solution_id>/proceed HTTP/1.1
Host: localhost:6971
```

### Organização do código

1. actions
    * São as principais ações do serviço, por exemplo, aprovar um reprocessamento, comitar a instância no domínio dentre outras
2. api
    * É a declaração da API do maestro;
3. broker
    * É onde está implementado a integração com o rabbitmq, este pacote contém as ações básicas de operação do broker;
4. etc
    * Pacote de funções utilitárias
5. handlers
    * Os handlers são os componentes que plugam no broker
6. models
    * Define o modelo de domínio usado pelo maestro
7. sdk
    * Implementa algumas chamadas de serviços da plataforma
8. vendor
    * É um pacote do Go onde ficam todas as bibliotecas de terceiros, os arquivos deste pacote jamais devem ser alterados diretamente;