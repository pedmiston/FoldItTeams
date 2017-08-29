import unipath

PROJ_ROOT = unipath.Path(__file__).absolute().ancestor(2)
PLAYBOOKS = PROJ_ROOT
R_PKG = PROJ_ROOT
R_PKG_DATA_RAW = R_PKG + '/data-raw'
R_PKG_DATA = R_PKG + '/data'
