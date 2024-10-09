# Service
    - [] Implement route
        - [x] Get All
        - [x] Create
        - [] Get Specific
        - [] Delete

# LoadBalancer 
    - [] Retry logic (route logic to another healthy instance when failed)    







# Note
    - Load balancer is working but it is currently forwarding the HTTP request to the service, we want to convert the http request to the plant
    order struct and use grpc call to the service with the struct