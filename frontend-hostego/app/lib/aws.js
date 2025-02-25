import AWS from "aws-sdk";

// Configure AWS SDK
AWS.config.update({
  region: process.env.NEXT_PUBLIC_AWS_REGION || "us-east-1",
  accessKeyId: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY,
  secretAccessKey: process.env.NEXT_PUBLIC_AWS_SECRET_KEY,
});

// Create S3 service object
const s3 = new AWS.S3();

export const uploadToS3Bucket = async (file) => {
  try {
    // Create unique file name using timestamp and random string
    const uniqueFileName = `hostego_img_${
      Date.now() * Math.random() * 12422532
    }.webp`;

    // Set up the upload parameters
    const params = {
      Bucket: "hostego-aws-bucket",
      Key: uniqueFileName,
      Body: file,
      ContentType: file.type,
    };

    // Using promise method from AWS SDK v2
    const { Location } = await s3.upload(params).promise();
    return Location;
  } catch (error) {
    console.error("S3 Upload Error:", error);
    throw new Error("Failed to upload file to S3");
  }
};
