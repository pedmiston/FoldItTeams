connect_to_db <- function() {
  DBI::dbConnect(RMySQL::MySQL(),
    host = "192.241.128.175",
    user = "foldit",
    password = Sys.getenv("FOLDIT_PASSWORD"),
    db = "Foldit"
  )
}
