# Sử dụng một hình ảnh base có sẵn của Node.js
FROM node:14

# Thiết lập thư mục làm việc trong container
WORKDIR /app

# Sao chép tất cả các tệp từ thư mục hiện tại của bạn vào /app trong container
COPY . .

# Cài đặt các dependency từ package.json
RUN npm install

# Mở cổng 3000 để giao tiếp với container này
EXPOSE 3000

# Khởi động ứng dụng khi container được chạy
CMD ["node", "index.js"]
