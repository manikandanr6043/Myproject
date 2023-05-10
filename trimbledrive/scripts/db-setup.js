// Create the following indexes in latest collection
// 1. Hashed index of spaceId to support sharding later.
// 2. Unique composite index of spaceId, parentFolderId, name. This will enforce uniqueness constraint on the combination of mentioned fields.
// 3. Non-Unique composite index of spaceId, parentFolderId, deleted. This will NOT enforce any constraints, these are used solely to improve query performance.
db.runCommand(
    {
        createIndexes: "latest",
        indexes: [
            {
                key: {"spaceId": "hashed"},
                collation: { locale: "en" },
                name: "sid"
            },
            {
                key: {
                    "spaceId": 1,
                    "parentFolderId": 1,
                    "name": 1
                },
                partialFilterExpression: {parentFolderId: {$exists: true}},
                name: "sid_pid_nm",
                unique: true,
            }, {
                key: {
                    "spaceId": 1,
                    "parentFolderId": 1,
                    "createdOnClock": -1,
                    "deleted": 1
                },
                collation: { locale: "en" },
                name: "sid_pid_cck_del"
            }, {
                key: {
                    "spaceId": 1,
                    "parentFolderId": 1,
                    "modifiedOnClock": -1,
                    "deleted": 1
                },
                collation: { locale: "en" },
                name: "sid_pid_mck_del"
            }
        ]
    }
)

// Create the following indexes in versions collection
// 1. Hashed index of spaceId to support sharding later.
// 2. Non-Unique composite index of spaceId, parentFolderId. This will NOT enforce any constraints, these are used solely to improve query performance.
db.runCommand(
    {
        createIndexes: "versions",
        indexes: [
            {
                key: {"spaceId": "hashed"},
                collation: { locale: "en" },
                name: "sid"
            },{
                key: {
                    "spaceId": 1,
                    "id": 1,
                    "modifiedOnClock": -1,
                    "createdOnClock": -1
                },
                collation: { locale: "en" },
                name: "sid_id_mck_cck"
            },{
                key: {
                    "spaceId": 1,
                    "id": 1,
                    "majorVersion": -1,
                    "minorVersion": -1
                },
                collation: { locale: "en" },
                name: "sid_id_mav_miv"
            }
        ]
    }
)
