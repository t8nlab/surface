import { path } from "@titanpl/native";
import { smtp } from "@titanpl/surface";

export default function sendmail(req) {
  try {
    const tplPath = path.resolve("../app/templates/contact.tmpl");

    const rawEmail = smtp.renderFile(tplPath, {
      name: "Soham",
      email: "ezetapp@gmail.com",
      from: "clashersoham07@gmail.com"
    });

    const result = smtp.send({
      host: "smtp.gmail.com",
      port: 587,
      username: "clashersoham07@gmail.com",
      password: "jjke wzkr tyfs aeod", 
      body: rawEmail,
      raw: true
    });

    return { success: true, result };
  } catch (err) {
    return { success: false, error: err.message };
  }
}
