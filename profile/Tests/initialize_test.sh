mongo --eval 'db.createUser({user: "testuser", pwd:"testpwd", roles:[{role:"userAdminAnyDatabase", db: "admin"}]})' admin
