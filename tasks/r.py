import invoke
from tasks import paths


@invoke.task
def install(ctx, collect_csvs_before_installing=False,
            use_data_before_installing=False):
    """Install the foldit R package."""
    if collect_csvs_before_installing:
        collect_csvs(ctx)
    if use_data_before_installing or collect_csvs_before_installing:
        use_data(ctx)
    ctx.run('cd {R_PKG} && Rscript install.R'.format(R_PKG=paths.R_PKG))


@invoke.task
def use_data(ctx, clear_data_before=False):
    """Compile the *.rda files to install with the R package."""
    if clear_data_before:
        ctx.run('rm -rf {R_PKG_DATA}'.format(R_PKG_DATA=paths.R_PKG_DATA),
                echo=True)
    ctx.run('cd {R_PKG} && Rscript data-raw/use-data.R'.format(R_PKG=paths.R_PKG))
