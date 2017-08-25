library(devtools)
library(readr)
library(dplyr)

TopScores <- read_csv("data-raw/top_scores.csv")
TopActions <- read_csv("data-raw/top_actions.csv")

use_data(
  TopScores,
  TopActions,
  overwrite = TRUE
)
