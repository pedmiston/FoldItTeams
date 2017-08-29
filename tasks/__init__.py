import glob
import json
import sqlite3
import invoke
import unipath
import pandas
import boto3

from tasks import s3
from tasks import r
from tasks import paths
from tasks.solution import TopSolution
from tasks.db import connect_to_db

namespace = invoke.Collection(s3, r)
