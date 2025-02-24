import { S3Client, PutObjectCommand } from "@aws-sdk/client-s3";

export const uploadToS3Bucket = async (file) => {
  console.log(process.env);
  console.log(file, "file is the new s3", process.env.NEXT_APP_ACCESS_KEY);

  const s3Client = new S3Client({
    region: "us-east-1", // specify your region
    credentials: {
      accessKeyId: "sdaf",
      secretAccessKey: "j",
    },
  });
  const params = {
    Bucket: "hostego-aws-bucket",
    Key: `hostego_img_${Date.now() * Math.random() * 12422532}.webp`,
    Body: file,
    ContentType: file.type,
  };

  try {
    const command = new PutObjectCommand(params);
    await s3Client.send(command);
    // Since PutObjectCommand doesn't return the URL directly, we need to construct it
    const image = `https://${params.Bucket}.s3.amazonaws.com/${params.Key}`;
    return image;
  } catch (error) {
    console.log(error);
  }
};
