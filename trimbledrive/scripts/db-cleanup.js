// Clean up all the databases created as part of the PR.

// Get all the db names from the cluster
var dbs = db.getMongo().getDBNames()

// Iterate through the db name list and delete the databases pertaining to the respective PR
for (var i in dbs) {
    if (dbs[i].includes(dbNamePrefix)) {
        db = db.getMongo().getDB(dbs[i]);
        db.dropDatabase();
    }
}
