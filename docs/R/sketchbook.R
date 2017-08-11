source("docs/R/setup.R")
# ---- sketchbook
label_group_type <- function(SolutionData) {
  SolutionData %>%
    mutate(group_type = ifelse(is.na(gid), "soloist", "evolver"))
}

recode_group_type <- function(SolutionData) {
  group_type_levels <- c("soloist", "evolver")
  group_type_labels <- c("Soloist", "Evolver")

  group_type_map <- data_frame(
    group_type = c("soloist", "evolver"),
    group_label = factor(group_type_levels, levels = group_type_levels,
                         labels = group_type_labels),
    group_c = c(-0.5, 0.5)
  )
  
  if (missing(SolutionData)) return(group_type_map)
  left_join(SolutionData, group_type_map)
}

label_outliers <- function(SolutionData) {
  SolutionData %>%
    mutate(is_outlier = ifelse(energy > 600, TRUE, FALSE))
}

data("TopSketchbook")

TopSketchbook <- TopSketchbook %>%
  label_group_type() %>%
  recode_group_type() %>%
  label_outliers()

set.seed(482)
energy_by_group_plot <- ggplot(TopSketchbook) +
  aes(group_label, energy) +
  geom_point(aes(shape = is_outlier),
             position = position_jitter(width = 0.1),
             alpha = 0.4) +
  geom_point(data = filter(TopSketchbook, is_outlier == FALSE),
             stat = "summary", fun.y = "mean", size = 3) +
  scale_y_reverse() +
  scale_shape_manual(values = c(1, 4), guide = "none")
