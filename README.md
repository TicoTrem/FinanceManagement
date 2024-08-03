# LifeManagement

Currently in development

This project is written in Go and keeps track of how much of your money isn't caught up in monthy expenses and savings contributions.
The console-based application includes an always running service, and a front end menu system. Both use a shared MySQL database. The structure of the project was designed
to make it easy for me to use TailScale to SSH in to my home server and launch an instance of the front end menu system. The console based nature also makes for fast development of the plethora of life management tools I plan to design.
Both the service and the front end menu system use a shared MySQL database.

- Users can input their expected monthly income, which will have monthly expenses and contributions to savings, and other financial goals, deducted from it to calculate that months 'Spending Money'.
- Users can input negative transactions when they spend money, lowering the spending money, and positive transactions when they make money, including their paychecks. Positive transactions do not instantly affect
spending money as those are already accounted for during the month through the expected monthly income.
- At the end of every month, the service will calculate your actual income that month, replacing the previous expected income in the spending money calculation.

This project was designed to provide a simpler alternative to something like the 'bucket' system of managing finances.
