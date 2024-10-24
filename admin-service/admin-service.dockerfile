# Wybór obrazu bazowego
FROM alpine:latest

# Tworzenie katalogu dla aplikacji
RUN mkdir /app

# Kopiowanie plików aplikacji
COPY userApp /app

# Uruchamianie aplikacji
CMD ["/app/userApp"]
