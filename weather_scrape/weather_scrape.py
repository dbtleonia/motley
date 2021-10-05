#!/usr/bin/python3

import bs4
import csv
import email
import sys

out = csv.writer(sys.stdout)
out.writerow(['year', 'time', 'temp'])

for year in range(2013, 2022):
    with open(f'{year}.mhtml') as f:
        for part in email.message_from_file(f).walk():
            if part.get_content_type() == 'text/html':
                soup = bs4.BeautifulSoup(part.get_payload(decode=True), features='html.parser')
                break

    for row in soup.find('lib-city-history-observation').find('tbody').find_all('tr'):
        cols = row.find_all('td')
        out.writerow([year, cols[0].get_text(), cols[1].get_text().split()[0]])
