# Сборка на основе пустого образа
FROM scratch

# Добавляем бинарный файл main в корень образа
ADD main /

# Здесь важно! Шаблоны и прочие внешние файлы в наш бинарник не включены
# Поэтому их следует скопировать в соответствующую папку контейнера
# Обратите внимание - при копировании пачки файлов по маске *.msg, каталог в контейнере должен заканчиваться закрывающим слэшем ...mail/

# Some ENV variables
# ENV PATH="/mm/mattermost/bin:${PATH}"

# Configure entrypoint and command
# COPY messaging-entrypoint.sh /entrypoint.sh

# RUN chmod +x /entrypoint.sh
# ENTRYPOINT ["/entrypoint.sh"]

# WORKDIR /mm/mattermost/bin
# CMD ["main", "server", "--config", "/mm/mattermost/config/config_docker.json"]

CMD ["main", "server", "--config", "config.json"]

EXPOSE 8065

