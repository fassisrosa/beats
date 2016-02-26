from openstackbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Openstackbeat normally
        """
        self.render_config_template(
                path=os.path.abspath(self.working_dir) + "/log/*"
        )

        openstackbeat_proc = self.start_beat()
        self.wait_until( lambda: self.log_contains("openstackbeat is running"))
        exit_code = openstackbeat_proc.kill_and_wait()
        assert exit_code == 0
