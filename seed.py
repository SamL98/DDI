import pymongo
import os

client = pymongo.MongoClient(os.environ['MONGODB_URI'])
db = client[os.environ['DBNAME']]

import pandas as pd

df = pd.read_csv('drugs.csv')
for i, row in df.iterrows():
	if i % 1000 == 0:
		print('%d / %d' % (i, len(df)))

	added = row['Drug_name1'].split(', ')
	base = row['Drug_name2'].split(', ')
	
	if len(base) == 1 and base[0] == 'BASELINE':
		base = []

	for drug in base:
		if drug in added:
			del added[added.index(drug)]

	assoc = {'added': added, 'base': base, 'or': row['OR'], 'p': row['Pvalue'], 'ci': row['95.CI'], 'ac': row['ac']}

	if len(base) == 0:
		db['S%d'%len(added)].insert_one(assoc)
	elif len(base) == 1:
		db['S%d'%(len(added)+3)].insert_one(assoc)
	else:
		db['S6'].insert_one(assoc)
