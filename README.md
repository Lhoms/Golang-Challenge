# Golang exercise  
 Important decisions/steps:  

- Price struct created. It is necessary to save creation date of that entry in the cache and check age. The cache will use this structure, but the service and wrapper interface will not change. Also if the cached price is expired the new one will be saved in the cache for the next time (with a new creation date). No changes in tests.

