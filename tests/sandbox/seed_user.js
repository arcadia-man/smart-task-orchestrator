const { MongoClient } = require('mongodb');
const bcrypt = require('bcrypt');

const uri = "mongodb://localhost:27017";
const client = new MongoClient(uri);

async function run() {
  try {
    await client.connect();
    const database = client.db("smart_orchestrator");
    const users = database.collection("users");

    // create hashed password
    const saltRounds = 10;
    const passwordHash = await bcrypt.hash("password123", saltRounds);

    const user = {
      name: "Harshil (Sandbox Tester)",
      email: "h@gmail.com",
      phone: "",
      password_hash: passwordHash,
      api_key: "68708e8e-aff3-4428-b1ce-2113ab247748",
      created_at: new Date(),
      updated_at: new Date(),
    };

    // Upsert the user
    await users.updateOne(
        { email: "h@gmail.com" },
        { $set: user },
        { upsert: true }
    );
    console.log("✅ User h@gmail.com seeded with API key 68708e8e-aff3-4428-b1ce-2113ab247748");
  } finally {
    await client.close();
  }
}

run().catch(console.dir);
