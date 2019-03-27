### reasons to have a seppareted service for the "scrapper"
  - Running with the "web" application when scaled could result in data duplicity if not treated correctly
  - Since it's a Small json being returned by the API I have decided to index all messages to a mongodb in case of API failure, in case of failure my  html page keeps working normally
  - There's no need to add all messages again, so i decided to create a checksum of the message and in case of content change, the new content will be added to my database as well.
  - Kubernetes Cronjob
  - 




TODO
- Prometheus exporter? in case of failure.
- Logs!




// Future
  (Scrapper plugin)
  Você registra um model e uma URL e ele busca usando go routines.


// Fazer feature toggle de count de documentos para que minha API não seja devastada e meu storage também.


// Deixar claro o motivo de não rodar um mongodb no kubernetes

// explicar o motivo desta arquitetura.