import { image } from "@titanpl/surface";

export default function thumbnail(req) {
  try {
    const userUploadedUrl = "https://i.pinimg.com/736x/2c/c2/fe/2cc2fe16eed28daf889d3fe5eff629c3.jpg";

    const result = image.resize({
      src: userUploadedUrl,
      width: 200,
      height: 0,
      quality: 100
    });

    return {
      success: true,
      message: "Processed natively without touching the disk",
      profilePicture: result.base64
    };
  } catch (err) {
    return {
      success: false,
      error: err.message
    };
  }
}
