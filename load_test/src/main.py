"""
Provide load test for gess node.
"""
from runner.environment import RunnerEnvironment
from runner import runner


if __name__ == '__main__':

    NODES_COUNT = 1
    LOAD_FACTOR = 1

    runner_environment = RunnerEnvironment(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)

    try:
        if runner_environment.setup_accounts():
            # runner.run(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)
            runner_environment.cleanup()

    except Exception as exception:
        runner_environment.save_accounts()

        raise exception
