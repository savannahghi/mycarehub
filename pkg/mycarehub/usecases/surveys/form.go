package surveys

import (
	"context"
	"fmt"
	"strings"
)

// getProgramIDs gets the program IDs for a given survey form
// These IDs are used to target programs that require unique set of surveys that are not shared with other programs;
// if program IDs are not are not included, the survey will be visible to all programs
// The programs object is defined as part of the instance data and the values are references from the setvalue section
// the ID values are separated by a full colon `:`
// below is a snippet of how it would look like in xForms.
/*
	<?xml version="1.0"?>
	<h:html xmlns="http://www.w3.org/2002/xforms" xmlns:ev="http://www.w3.org/2001/xml-events" xmlns:h="http://www.w3.org/1999/xhtml" xmlns:jr="http://openrosa.org/javarosa" xmlns:odk="http://www.opendatakit.org/xforms" xmlns:orx="http://openrosa.org/xforms" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
		<h:head>
			<model odk:xforms-version="1.0.0">
				<instance>
					<data id="akmCQQxf4LaFjAWDbg29pj" version="V0.0.9">
						<programs/>
					</data>
				</instance>
				<setvalue event="odk-instance-first-load" ref="/data/programs" value="'4181df12-ca96-4f28-b78b-8e8ad88b25df:5181df12-ca96-4f28-b78b-8e8ad88b25df'"/>
			</model>
		</h:head>
	</h:html>
*/
func getProgramIDs(ctx context.Context, form map[string]interface{}) ([]string, error) {

	programIDs := []string{}

	formHTML, ok := form["html"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid form: expected a 'html' key")
	}

	formHead, ok := formHTML["head"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid form: expected a 'head' key")
	}

	formModel, ok := formHead["model"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid form: expected a 'model' key")
	}

	// skip surveys that do not contain program IDs
	formSetValue, ok := formModel["setvalue"].(map[string]interface{})
	if !ok {
		return programIDs, nil
	}

	formReference, ok := formSetValue["-ref"].(string)
	if !ok {
		return programIDs, nil
	}

	if formReference != "/data/programs" {
		return programIDs, nil
	}

	programValue, ok := formSetValue["-value"].(string)
	if !ok {
		return programIDs, nil
	}

	programValue = strings.ReplaceAll(programValue, "'", "")

	programIDs = strings.Split(programValue, ":")

	return programIDs, nil
}
