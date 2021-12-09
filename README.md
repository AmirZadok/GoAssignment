# GoAssignment
  REST API to manage containers in a database.<br>
  The server call examples on localhost : 
  - 
  <ul>
  <li>Get all hosts list  -  http://localhost:80/host</li>
  <li>Get all containers list - http://localhost:80/container</li>
  <li>Get Host by It - http://localhost:80/host/(hostId)</li>
  <li>Get containers by ID - http://localhost:80/container/{containerId}</li>
  <li>Get all containers for specific host - http://localhost:80/container/container-for-spec-host/{hostId} . also check there is such hostId in hosts.</li>
  <li>Insert container - http://localhost:80/container . receive the params from the POST body as json. generate into container.name uuid</li>
   
</ul>
  
  
