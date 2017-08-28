import glob
import json
import sqlite3
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
def push_to_db(ctx):
    """Collect data from json files into csvs for analysis."""
    con = foldit.connect_to_db()

    for json_file in glob.glob(foldit.TOP_SOLUTION_DATA_GLOB):
        print("Processing json file: " + json_file)
        top_scores = []
        top_actions = []

        for solution_data in open(json_file):
            top_solution = foldit.TopSolution(json.loads(solution_data))
            top_scores.append(top_solution.to_record())
            top_actions.append(top_solution.get_actions())
        
        pandas.DataFrame.from_records(top_scores).to_sql(
            'TopScores', con, if_exists='append')
        pandas.concat(top_actions, ignore_index=True).to_sql(
            'TopActions', con, if_exists='append')

