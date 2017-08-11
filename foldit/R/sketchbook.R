#' Top solutions to sketchbook puzzles with a fixed number of moves.
#'
#' \describe{
#' \item{uid}{factor Unique identifier for users.}
#' \item{gid}{factor Unique identifier for groups. NA for soloists.}
#' \item{timestamp}{int Time at which this solution was submitted.}
#' \item{energy}{double Energy score of winning protein configuration.}
#' \item{puzzle_id}{factor Unique identifier for puzzles.}
#' \item{prefix}{char Type of rank comparison.}
#' \item{rank}{int Ranked energy score by prefix.}
#' }
"TopSketchbook"
