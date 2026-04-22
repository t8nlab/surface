import t from "@titanpl/route";

t.get("/u").action("getuser")
t.get("/m").action("sendmail")
t.get("/mb").action("bulkmail")
t.get("/ir").action("thumbnail")
t.get("/json").action("testjson")
t.get("/cloud").action("cloud_stream")
t.get("/clean").action("testclean")
t.get("/extract").action("testextract")


t.get("/").reply("Ready to land on Titan Planet 🚀");

t.start(5100, "Titan Running!", 12);
