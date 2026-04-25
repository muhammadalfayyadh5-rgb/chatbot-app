# Menggunakan image Nginx yang ringan
FROM nginx:alpine

# Menyalin seluruh isi folder saat ini ke folder default Nginx di dalam container
COPY . /usr/share/nginx/html

# Expose port 80 agar bisa diakses
EXPOSE 80

# Jalankan Nginx
CMD ["nginx", "-g", "daemon off;"]