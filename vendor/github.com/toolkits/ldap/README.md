[![GoDoc](https://godoc.org/github.com/vanackere/ldap?status.svg)](https://godoc.org/github.com/vanackere/ldap) [![Build Status](https://travis-ci.org/vanackere/ldap.svg)](https://travis-ci.org/vanackere/ldap)

Basic LDAP v3 functionality for the GO programming language.
------------------------------------------------------------

* Required library:
 - github.com/vanackere/asn1-ber

* Working:
 - Connecting to LDAP server
 - Binding to LDAP server
 - Searching for entries
 - Compiling string filters to LDAP filters
 - Paging Search Results
 - Modify Requests / Responses

* Examples:
 - search
 - modify

* Tests Implemented:
-   Filter Compile / Decompile

* TODO:
-   Add Requests / Responses
-   Delete Requests / Responses
-   Modify DN Requests / Responses
-   Compare Requests / Responses
-   Implement Tests / Benchmarks


This feature is disabled at the moment, because in some cases the "Search Request Done" packet will be handled before the last "Search Request Entry":
 -   Mulitple internal goroutines to handle network traffic
      Makes library goroutine safe
      Can perform multiple search requests at the same time and return
         the results to the proper goroutine.  All requests are blocking
         requests, so the goroutine does not need special handling
