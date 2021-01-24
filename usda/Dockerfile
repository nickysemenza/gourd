FROM debian:buster

RUN apt-get update && apt-get install -y unzip postgresql

ADD https://fdc.nal.usda.gov/fdc-datasets/FoodData_Central_csv_2020-04-29.zip /tmp
# ADD  FoodData_Central_csv_2020-04-29.zip /tmp
WORKDIR /work
COPY import.sh .
RUN unzip /tmp/FoodData_Central_csv_2020-04-29.zip -d csv/
RUN chmod +x /work/import.sh
ENTRYPOINT ["bash","/work/import.sh", "/work/csv/"]