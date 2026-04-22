import axios from "axios"

const minio = axios.create({
  timeout: 30000,
});

export default minio