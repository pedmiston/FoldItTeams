import glob
import invoke
import unipath
from scripts.filter_already_downloaded import SOLUTION_DATA

PROJ_ROOT = unipath.Path(__file__).absolute().parent
R_PKG = PROJ_ROOT


@invoke.task
def install(ctx, use_data_before_installing=False):
    """Install the foldit R package."""
    if use_data_before_installing:
        use_data(ctx)
    ctx.run('cd {R_PKG} && Rscript install.R'.format(R_PKG=R_PKG))


@invoke.task
def use_data(ctx, clear_data_before=False):
    """Compile the *.rda files to install with the R package."""
    if clear_data_before:
        ctx.run('rm -rf {R_PKG_DATA}'.format(R_PKG_DATA=unipath.Path(R_PKG, 'data')),
                echo=True)
    ctx.run('cd {R_PKG} && Rscript data-raw/use-data.R'.format(R_PKG=R_PKG))


@invoke.task
def collect_csvs(ctx):
    """Collect data from json files into csvs for analysis."""
    for json_file in glob.glob(SOLUTION_DATA):
        pass
