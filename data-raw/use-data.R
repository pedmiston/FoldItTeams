library(devtools)
library(readr)
library(dplyr)

TopSketchbook <- read_csv("../scripts/data/sketchbook_top.csv") %>%
  mutate(
    uid = factor(uid),
    gid = factor(ifelse(gid != 0, gid, NA))
  )

use_data(TopSketchbook, overwrite = TRUE)
