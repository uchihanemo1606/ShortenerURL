URL shortener là dịch vụ rút gọn URL dài thành đoạn mã ngắn, và theo dõi được số lượt click của shortURL 

## Mô tả bài toán
dịch vụ này giải quyết được các vấn đề sau:

1. Rút ngắn URL dài thành mã ngắn (ví dụ : từ https://github.com/upstash/context7/blob/f2f367d8913843bd28b2a96a6ce860f43e3fc3ca/Dockerfile thành localhost:8080/abc123)
2. Khi truy cập URL ngắn thì tự động chuyển đến URL gốc
3. theo dõi lượt click khi có người truy cập URL ngắn
4. Xác thực: chỉ những người đã login mới có thể tạo URL ngắn
5. mỗi URL ngắn có thời gian sử dụng tối đa 14 ngày
6. validate URL trước khi tạo URL ngắn (không cho phép tạo nhiều URL ngắn từ một URL gốc)

## Yêu cầu hệ thống
- Golang phiên bản 1.25.0
- databse: redis(nên sử dụng cloud từ Upstash Redis để đơn giản)


## Cài đặt project
1. clone project về máy: git clone https://github.com/uchihanemo1606/ShortenerURL
2. di chuyển đến thư mục chứa project cd URLShortener

## Cài đặt dependencies
go mod download

## Cài đặt biến môi trường env
1. Tạo file .env ở vị trí gốc của project(cùng cấp với file go.mod, và go.sum)
2. Thêm giá trị đúng cho REDIS_URL = (giá trị này được cung cấp ở trong database phần detail REDIS_URL có format như sau: rediss://default:yourpassword@literate-egret.upstash.io)
3. Thêm giá trị cho JWT_SECRET_KEY = (giá trị này tuỳ ý nên để ngẫu nhiên để tăng tính bảo mật)
4. Thêm giá trị cho BASE_URL = đây là host và port mà backend chạy(ví dụ: localhost:8080 nếu là server local)

## Chạy redis(nếu không dùng redis cloud)
1. Cài đặt Redis (Windows)
2. Download từ https://redis.io/download và cài đặt
3. Chạy Redis server

## Chạy ứng dụng
1. go mod tidy
2. go run main.go(hoặc chỉ cần gõ "air" trên terminal)

## Test API
Server sẽ chạy tại `http://localhost:8080`
các enpoint có sẵn:
Các endpoint có sẵn:
   - `POST /signup` - Đăng ký tài khoản
   - `POST /login` - Đăng nhập (sau khi login xong thì copy lại token được cung cấp)
   - `POST /shorten?url=<long_url>` - Tạo URL ngắn (past token đã copy vào Authorization->AuthType(Bearer Token))
   - `GET /<short_code>` - Chuyển hướng đến URL gốc
   - `GET /urls` - Lấy danh sách tất cả URLs

1. Đăng ký user:
   ```bash
   curl -X POST http://localhost:8080/signup \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. Đăng nhập:
   ```bash
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```
   Sẽ trả về JWT token.

3. Tạo URL ngắn:
   ```bash
   curl -X POST "http://localhost:8080/shorten?url=https://google.com" \
     -H "Authorization: Bearer <your-jwt-token>"
   ```

4. Truy cập URL ngắn:
   ```bash
   curl -L http://localhost:8080/<short-code>
   ```

## Thiết kế và chọn kỹ thuật

1. lý do chọn redis làm database:
   - Redis là 1 database noSQL dễ cài dặt và đọc ghi rất nhanh, rất phù hợp với URL shortener không yêu cầu cao đối với thời gian phản hồi.
   - Rất dễ cài đặt không yêu cầu cao về kỹ thuật, không cần schema, lưu trực tiếp bằng JSON luôn.
   - Có thể set expiration cho key khi URL hết hạn thì nó sẽ bị redis xoá(phải setting trên redis upstash trước)
   - Không cần lưu nhiều dữ liệu nên phiên bản free tier có thể lưu được hàng triệu URL
  
2. Lý do dùng RESTful API
   - Đơn giản và rất phổ biến: RESTful API sử dụng methods GET/POST đơn giản, rất nhiều người dùng để phát triển service
   - Không tốn nhiều thời gian học tập đối với người mới
     
3. Giải thích thuật toán generate mã ngắn và giải quyết conflict
   - Validate URL gốc: kiểm tra tính hợp lệ của URL trước khi tạo mã ngắn
   - Kiểm tra trùng lặp: dùng URL gốc để query đến redis nếu có kết quả trả về tức là URL gốc đã được dùng để tạo URL ngắn trả kết quả và dừng function ngay lập tức
   - dùng thư viện `crypto/rand` để tạo mã ngắn ngẫu nhiên gồm 6 ký tự
   - Kiểm tra trùng lặp của mã ngắn: dùng mã ngắn để query đến redis nếu có kết quả tức là mã ngắn đã được dùng để định danh URL gốc quay lại phần tạo mã ngắn và kiểm tra cho đến khi không có mã ngắn trùng lặp
   - dùng userID lấy từ context để quản lý URL ngắn do ai tạo
   - lưu vào redis và kết thúc hàm

## Trade-offs
1. lý do chọn redis thay vì Mongodb hay các CSDL SQL
  ###ưu điểm
   - hiệu năng: nếu có 10-100x query thì redis cho tốc độ nhanh hơn
   - tính đơn giản: không cần thiết kế database phức tạp, nếu có lỗi thì không cần đập đi xây lại
   - Tiêu hao ít tài nguyên hơn
  ###nhược điểm
  - Không có complex queries (JOIN, aggregation)
  - dữ liệu không có cấu trúc khi thay đổi cấu trúc thì rất khó quản lý
=> phù hợp với dự án vì nó chủ yếu key-value không cần chính xác tuyệt đối với dữ liệu
2. Chọn random generation thay vì sequential IDs
  ###Ưu điểm
   - Không đoán được URL tiếp theo (tránh việc user dự đoán url)
   - Không expose thông tin về số lượng URLs
   - đơn giản
  ###nhược điểm
   - có thể bị trùng mã ngắn(khả năng thấp nếu ít url)
   - khó test hàng loạt
=> phù hợp vì: tính đơn giản và bảo mật

## Challenges
Challenge 1: Xử lý đồng thời nhiều requests cho cùng URL
hướng giải quyết : dùng lock chỉ giải quyết 1 request tạo URL trên cùng 1 thời điểm
học được: tác hại và sự nguy hiểm của deadlock đối với dữ án

Challenge 2: Đảm bảo tính duy nhất đối với random generation
hướng giải quyết: dùng vòng lặp để random tạo mã ngắn cho đến khi không trùng lặp(nếu việc trùng lặp xảy ra liên tục trên nhiều request thì sẽ gây ra chậm, lag server)
học được: nhìn nhận được tầm quan trọng của việc kiểm soát lỗi các lỗi nhỏ có thế gây hư hỏng hệ thống

Challenge 3: nhập URL không hợp lệ nhưng hệ thống vẫn cho ra kết quả
hướng giải quyết: dùng thư viện `http/net` của golang để validate URL gốc 
học được: tầm quan trọng của việc validate dữ liệu

##Limitations & Improvements
Code hiện tại còn thiếu: 
  - Metrics: chưa có heal checks
  - Rate Limiting: chưa giới hạn số lượng requests của user
  - Input Validation: validate vẫn còn rất cơ bản có thể bị hacker lách
  - Documentation: chưa có tài liệu API
  - Logging: còn rất sơ sài

## Nếu có thêm thời gian em sẽ làm:
  - viết tài liệu API
  - thêm bảo mật cho hệ thống: HTTPS,CORS,CSRF
  - cải thiện loggin
  - hạn chế số lượng request trên cùng 1 thời điểm
  - Validate dữ liệu kĩ hơn

## Production-ready cần thêm
  - sercurity; SSL/TLS
  - backup: cần backup dữ liệu database và code
  - Test: Test Bug trước tri lên production
  - Cải thiện hiệu năng: cải thiện code để năng cao hiệu năng
  - viết DockerFile
