import glob
import json
import invoke
import unipath
import pandas
import feather
import foldit


@invoke.task
def install(ctx, collect_csvs_before_installing=False,
            use_data_before_installing=False):
    """Install the foldit R package."""
    if collect_csvs_before_installing:
        collect_csvs(ctx)
    if use_data_before_installing or collect_csvs_before_installing:
        use_data(ctx)
    ctx.run('cd {R_PKG} && Rscript install.R'.format(R_PKG=foldit.R_PKG))


@invoke.task
def use_data(ctx, clear_data_before=False):
    """Compile the *.rda files to install with the R package."""
    if clear_data_before:
        ctx.run('rm -rf {R_PKG_DATA}'.format(R_PKG_DATA=foldit.R_PKG_DATA),
                echo=True)
    ctx.run('cd {R_PKG} && Rscript data-raw/use-data.R'.format(R_PKG=foldit.R_PKG))


@invoke.task
def collect_csvs(ctx):
    """Collect data from json files into csvs for analysis."""
    top_scores_dst = foldit.R_PKG_DATA_RAW + '/top_scores.csv'
    top_actions_dst = foldit.R_PKG_DATA_RAW + '/top_actions.csv'

    top_scores = []
    top_actions = []

    for json_file in glob.glob(foldit.TOP_SOLUTION_DATA_GLOB):
        for solution_data in open(json_file):
            top_solution = foldit.TopSolution(json.loads(solution_data))
            top_scores.append(top_solution.to_record())
            top_actions.append(top_solution.get_actions())

    pandas.DataFrame.from_records(top_scores).to_csv(top_scores_dst, index=False)
    pandas.concat(top_actions).to_csv(top_actions_dst, index=False)

@invoke.task
def collect_feather(ctx):
    """Collect data from json files into dataframes in feather format."""
    top_scores_dst = foldit.R_PKG_DATA_RAW + '/top_scores.feather'
    top_actions_dst = foldit.R_PKG_DATA_RAW + '/top_actions.feather'

    top_scores = []
    top_actions = []

    for json_file in glob.glob(foldit.TOP_SOLUTION_DATA_GLOB):
        for solution_data in open(json_file):
            top_solution = foldit.TopSolution(json.loads(solution_data))
            top_scores.append(top_solution.to_record())
            top_actions.append(top_solution.get_actions())

    feather.write_dataframe(pandas.DataFrame.from_records(top_scores),
                            top_scores_dst)
    feather.write_dataframe(pandas.concat(top_actions),
                            top_actions_dst)
