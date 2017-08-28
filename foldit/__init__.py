import unipath
from foldit.solution import TopSolution
from foldit.db import connect_to_db

PROJ_ROOT = unipath.Path(__file__).absolute().ancestor(2)
PLAYBOOKS = PROJ_ROOT
R_PKG = PROJ_ROOT
R_PKG_DATA_RAW = R_PKG + '/data-raw'
TOP_DATA_DIR = R_PKG_DATA_RAW + '/top'
TOP_SOLUTION_DATA_GLOB = TOP_DATA_DIR + '/run-*/solution_data.json'

R_PKG_DATA = R_PKG + '/data'
