# Golang exercise  
 Important decisions/steps:  

- Price struct created. It is necessary to save creation date of that entry in the cache and check age. The cache will use this structure, but the service and wrapper interface will not change. Also if the cached price is expired the new one will be saved in the cache for the next time (with a new creation date). No changes in tests.

- To parallelize price request I used go routines and a channel to retrieve the response.

- To communicate in the channel I use the PriceRequest Struct, with this I can get the price(number) or the error.

- All changes were made in GetPricesFor and not in GetPriceFor, don't need to be affected with this new paralleled implementation. 

- Test added. The existing test was not checking GetPricesFor in case of error, only GetPriceFor.
