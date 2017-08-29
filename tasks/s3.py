import invoke
import boto3
import botocore
from tasks import paths


@invoke.task
def push(ctx):
    """Push data to S3 bucket."""
    filename = paths.R_PKG_DATA_RAW + '/top.json'
    assert filename.exists()
    s3 = boto3.resource('s3')
    data = open(filename, 'rb')
    s3.Bucket('foldit').put_object(Key='top.json', Body=data)
    data.close()


@invoke.task
def pull(ctx):
    """Pull data from S3 bucket."""
    filename = paths.R_PKG_DATA_RAW + '/top.json'
    s3 = boto3.resource('s3')
    try:
        s3.Bucket('foldit').download_file('top.json', filename)
    except botocore.exceptions.ClientError as e:
        if e.response['Error']['Code'] == "404":
            print("The object does not exist")
        else:
            raise
