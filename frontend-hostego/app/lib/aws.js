import AWS from "aws-sdk";

// Only log in development
if (process.env.NODE_ENV !== "production") {
  console.log("AWS Environment Variables:", {
    region: process.env.NEXT_PUBLIC_AWS_REGION,
    accessKey: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY,
    secretKey: process.env.NEXT_PUBLIC_AWS_SECRET_KEY,
    bucket: process.env.NEXT_PUBLIC_AWS_BUCKET_NAME,
  });
}

// Configure AWS SDK
const awsConfig = {
  region: process.env.NEXT_PUBLIC_AWS_REGION,
  credentials: {
    accessKeyId: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY,
    secretAccessKey: process.env.NEXT_PUBLIC_AWS_SECRET_KEY,
  },
};

// Initialize AWS with configuration
AWS.config.update(awsConfig);

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
      Bucket: process.env.NEXT_PUBLIC_AWS_BUCKET_NAME,
      Key: uniqueFileName,
      Body: file,
      ContentType: file.type,
    };

    // Using promise method from AWS SDK v2
    const { Location } = await s3.upload(params).promise();
    return Location;
  } catch (error) {
    // More detailed error logging in development
    if (process.env.NODE_ENV !== "production") {
      console.error("S3 Upload Error:", error);
    }
    throw new Error("Failed to upload file to S3");
  }
};
