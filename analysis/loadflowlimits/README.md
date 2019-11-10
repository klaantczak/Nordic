Calculates initial flow through the nordic32 network links.

To set the initial flow through the links, add the following code to the
loadflow plugin:

 func (p *Plugin) Init(e hps.IEnvironment) error {
 	substations, _ := q.FindSubstationsNetwork(e)
 	model := BuildModel(substations)
 
+	states := buildStates(model)
+
+	result := loadflow(model, states)
+
+	for i, v := range result.Flows {
+		machine := model.Links[i].Machine
+		log.Printf("flow['%v'] = %v", machine.Name(), v)
+	}
+ 
 	model.Modified = false

Build the project, run "./prepare.sh", copy the data from result.jslog,
update the numbers in "update.js" and run "update.js".