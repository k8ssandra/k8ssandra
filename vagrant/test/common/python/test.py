from cassandra.cluster import Cluster
from cassandra.auth import PlainTextAuthProvider

auth_provider = PlainTextAuthProvider(username='k8ssandra-superuser', password='UBoPcvdrm-u8jw0Tz3bJZTTeuRhHoiNTwtVpjJ9-EVh0ePhYzDhtgA')
cluster = Cluster(['localhost'], port=9042, auth_provider=auth_provider)
session = cluster.connect()
session.set_keyspace('test_rest')

session.execute("""
        CREATE TABLE IF NOT EXISTS test_table (
            testkey text,
            testcol1 text,
            testcol2 text,
            PRIMARY KEY (testkey, testcol1)
        )
        """)

rows = session.execute("""
        SELECT * FROM test_table
        """)

print rows
