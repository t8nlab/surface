import t from "@titanpl/route";

t.get("/u").action("getuser")
t.get("/m").action("sendmail")
t.get("/mb").action("bulkmail")
t.get("/ir").action("thumbnail")

t.get("/").reply("Ready to land on Titan Planet 🚀");

t.start(5100, "Titan Running!");
