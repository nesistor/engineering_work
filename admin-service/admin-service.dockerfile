# Wybór obrazu bazowego
FROM alpine:latest

# Tworzenie katalogu dla aplikacji
RUN mkdir /app

# Kopiowanie plików aplikacji
COPY adminApp /app

# Uruchamiaanie aplikacji
CMD ["/adminApp"]
